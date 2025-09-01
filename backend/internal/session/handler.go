package session

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/studyplatform/backend/pkg/database"
	"github.com/studyplatform/backend/pkg/models"
)

// StartSessionHandler starts a new study session
func StartSessionHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
			return
		}
		var req models.StartSessionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		sessions := mongoClient.GetCollection(database.CollectionNames.Sessions)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		sessionMap := bson.M{
			"user_id":            userIDStr,
			"room_id":            req.RoomID,
			"start_time":         time.Now(),
			"is_active":          true,
			"xp_earned":          0,
			"pomodoro_completed": 0,
			"inactivity_periods": []bson.M{},
		}
		res, err := sessions.InsertOne(ctx, sessionMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start session"})
			return
		}
		var session models.Session
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			session.ID = oid
		}
		err = sessions.FindOne(ctx, bson.M{"_id": session.ID}).Decode(&session)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created session"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"session": session.ToResponse()})
	}
}

// EndSessionHandler ends a study session
func EndSessionHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
			return
		}

		var req struct {
			SessionID        string            `json:"sessionId" binding:"required"`
			Duration         int64             `json:"duration" binding:"required"` // In seconds
			InactiveDuration int64             `json:"inactiveDuration"`            // In seconds
			PomodoroCount    int               `json:"pomodoroCount"`               // Number of completed pomodoros
			ActivityData     []models.Activity `json:"activityData"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		sessionObjID, err := primitive.ObjectIDFromHex(req.SessionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
			return
		}

		sessions := mongoClient.GetCollection(database.CollectionNames.Sessions)
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Find the session
		var session models.Session
		err = sessions.FindOne(ctx, bson.M{"_id": sessionObjID, "user_id": userIDStr, "is_active": true}).Decode(&session)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Active session not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Calculate XP earned
		xpEarned := calculateXPEarned(req.Duration, req.InactiveDuration, req.PomodoroCount)

		// Update session
		endTime := time.Now()
		sessionUpdate := bson.M{
			"$set": bson.M{
				"end_time":           endTime,
				"duration":           req.Duration,
				"inactive_time":      req.InactiveDuration,
				"xp_earned":          xpEarned,
				"is_active":          false,
				"pomodoro_completed": req.PomodoroCount,
				"activities":         req.ActivityData,
				"updated_at":         time.Now(),
			},
		}

		_, err = sessions.UpdateOne(ctx, bson.M{"_id": sessionObjID}, sessionUpdate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
			return
		}

		// Update user XP
		userUpdate := bson.M{
			"$inc": bson.M{"xp": xpEarned},
			"$set": bson.M{"last_active": time.Now()},
		}

		_, err = users.UpdateOne(ctx, bson.M{"unique_id": userIDStr}, userUpdate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user XP"})
			return
		}

		// Get updated session
		err = sessions.FindOne(ctx, bson.M{"_id": sessionObjID}).Decode(&session)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated session"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"session":  session.ToResponse(),
			"xpEarned": xpEarned,
			"message":  "Session ended successfully",
		})
	}
}

// ActivityPingHandler updates user activity to prevent session timeout
func ActivityPingHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
			return
		}

		sessionID := c.Param("id")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Session ID required"})
			return
		}

		sessionObjID, err := primitive.ObjectIDFromHex(sessionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
			return
		}

		var req struct {
			ActivityType string `json:"activityType" binding:"required"`
			Timestamp    string `json:"timestamp"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		sessions := mongoClient.GetCollection(database.CollectionNames.Sessions)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Create activity record
		activity := models.Activity{
			Type:      req.ActivityType,
			Timestamp: time.Now(),
			Details:   "User activity ping",
		}

		// Update session with new activity
		update := bson.M{
			"$push": bson.M{"activities": activity},
			"$set":  bson.M{"updated_at": time.Now()},
		}

		result, err := sessions.UpdateOne(ctx, bson.M{"_id": sessionObjID, "user_id": userIDStr, "is_active": true}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Active session not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Activity recorded"})
	}
}

// GetUserSessionStats returns comprehensive session statistics for a user
func GetUserSessionStats(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
			return
		}

		sessions := mongoClient.GetCollection(database.CollectionNames.Sessions)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Aggregate session statistics
		pipeline := []bson.M{
			{"$match": bson.M{"user_id": userIDStr, "is_active": false}},
			{"$group": bson.M{
				"_id":                nil,
				"totalSessions":      bson.M{"$sum": 1},
				"totalDuration":      bson.M{"$sum": "$duration"},
				"totalXPEarned":      bson.M{"$sum": "$xp_earned"},
				"totalInactiveTime":  bson.M{"$sum": "$inactive_time"},
				"totalPomodoroCount": bson.M{"$sum": "$pomodoro_completed"},
				"averageDuration":    bson.M{"$avg": "$duration"},
			}},
		}

		cursor, err := sessions.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate statistics"})
			return
		}
		defer cursor.Close(ctx)

		var stats []bson.M
		if err = cursor.All(ctx, &stats); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode statistics"})
			return
		}

		if len(stats) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"totalSessions":      0,
				"totalDuration":      0,
				"totalXPEarned":      0,
				"totalInactiveTime":  0,
				"totalPomodoroCount": 0,
				"averageDuration":    0,
			})
			return
		}

		c.JSON(http.StatusOK, stats[0])
	}
}

// CheckXPPrivileges checks if a user has sufficient XP for certain actions
func CheckXPPrivileges(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
			return
		}

		action := c.Query("action")
		if action == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Action parameter required"})
			return
		}

		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.User
		err := users.FindOne(ctx, bson.M{"unique_id": userIDStr}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			return
		}

		privileges := calculateUserPrivileges(user.XP)

		switch action {
		case "create_shared_room":
			c.JSON(http.StatusOK, gin.H{
				"hasPrivilege": privileges.CanCreateSharedRoom,
				"required":     1000,
				"current":      user.XP,
				"remaining":    max(0, 1000-user.XP),
			})
		case "add_participant":
			maxParticipants := privileges.MaxParticipants
			c.JSON(http.StatusOK, gin.H{
				"hasPrivilege":    true,
				"maxParticipants": maxParticipants,
				"current":         user.XP,
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown action"})
		}
	}
}

// XPPrivileges represents user privileges based on XP
type XPPrivileges struct {
	CanCreateSharedRoom bool
	MaxParticipants     int
	MaxSharedRooms      int
}

// calculateUserPrivileges determines user privileges based on XP
func calculateUserPrivileges(xp int) XPPrivileges {
	privileges := XPPrivileges{
		CanCreateSharedRoom: xp >= 1000,
		MaxParticipants:     4, // Default: 4 participants (5 total including owner)
		MaxSharedRooms:      1, // Default: 1 shared room
	}

	// Each additional 300 XP allows +1 participant (max 10 total)
	if xp >= 300 {
		additionalParticipants := (xp - 300) / 300
		privileges.MaxParticipants = min(4+additionalParticipants, 9) // Max 9 additional (10 total)
	}

	// Each additional 1000 XP allows +1 shared room (max 5 total)
	if xp >= 2000 {
		additionalRooms := (xp - 1000) / 1000
		privileges.MaxSharedRooms = min(1+additionalRooms, 5)
	}

	return privileges
}

// calculateXPEarned calculates XP based on session data
func calculateXPEarned(duration, inactiveDuration int64, pomodoroCount int) int {
	// Base XP calculation: 2 XP per minute of active time
	activeDuration := duration - inactiveDuration
	activeMinutes := activeDuration / 60
	baseXP := int(activeMinutes * 2)

	// Bonus XP for completed pomodoros: 30 XP per pomodoro
	pomodoroXP := pomodoroCount * 30

	// Inactivity penalty: lose 1 XP per 5 minutes of inactivity
	inactiveMinutes := inactiveDuration / 60
	penalty := int(inactiveMinutes / 5)

	totalXP := baseXP + pomodoroXP - penalty

	// Minimum XP is 0
	if totalXP < 0 {
		totalXP = 0
	}

	return totalXP
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ListSessionsHandler returns all sessions for the authenticated user
func ListSessionsHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userIDStr, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
			return
		}

		sessions := mongoClient.GetCollection(database.CollectionNames.Sessions)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{"user_id": userIDStr}
		cursor, err := sessions.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer cursor.Close(ctx)

		var sessionList []models.SessionForResponse
		for cursor.Next(ctx) {
			var session models.Session
			if err := cursor.Decode(&session); err == nil {
				sessionList = append(sessionList, session.ToResponse())
			}
		}
		c.JSON(http.StatusOK, gin.H{"sessions": sessionList})
	}
}
