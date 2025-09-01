package room

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/studyplatform/backend/pkg/database"
	"github.com/studyplatform/backend/pkg/models"
)

// ListRoomsHandler returns all rooms for the authenticated user
func ListRoomsHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// First, try to find the user to get their unique_id
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		var user models.User
		userFilter := bson.M{"$or": []bson.M{
			{"_id": userIDStr},       // Try by MongoDB ObjectID
			{"unique_id": userIDStr}, // Try by unique_id
		}}
		err := users.FindOne(ctx, userFilter).Decode(&user)
		if err != nil {
			// If user not found, just use the original userIDStr
			user.UniqueID = userIDStr
		}

		// Filter rooms where user is creator or participant
		filter := bson.M{"$or": []bson.M{
			{"creator_id": user.UniqueID},   // User is creator (creator_id is stored as string)
			{"participants": user.UniqueID}, // User is participant (participants array contains strings)
		}}

		cursor, err := rooms.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer cursor.Close(ctx)

		var roomList []models.RoomForResponse
		for cursor.Next(ctx) {
			var room models.Room
			if err := cursor.Decode(&room); err == nil {
				// Get creator's username
				var creator models.User
				creatorFilter := bson.M{"unique_id": room.CreatorID}
				if err := users.FindOne(ctx, creatorFilter).Decode(&creator); err == nil {
					roomResponse := room.ToResponse()
					roomResponse.CreatorUsername = creator.Username

					// Populate participant details
					var participantInfos []models.ParticipantInfo
					for _, participantID := range room.Participants {
						var participant models.User
						participantFilter := bson.M{"unique_id": participantID}
						if err := users.FindOne(ctx, participantFilter).Decode(&participant); err == nil {
							participantInfos = append(participantInfos, models.ParticipantInfo{
								UserID:    participant.UniqueID,
								Username:  participant.Username,
								AvatarURL: participant.AvatarURL,
								IsOnline:  participant.IsActive, // Use IsActive as a simple online indicator
							})
						} else {
							// If participant not found, still add them with basic info
							participantInfos = append(participantInfos, models.ParticipantInfo{
								UserID:    participantID,
								Username:  "Unknown User",
								AvatarURL: "",
								IsOnline:  false,
							})
						}
					}
					roomResponse.Participants = participantInfos

					roomList = append(roomList, roomResponse)
				} else {
					// If creator not found, still add room but with empty username
					roomResponse := room.ToResponse()

					// Populate participant details even if creator not found
					var participantInfos []models.ParticipantInfo
					for _, participantID := range room.Participants {
						var participant models.User
						participantFilter := bson.M{"unique_id": participantID}
						if err := users.FindOne(ctx, participantFilter).Decode(&participant); err == nil {
							participantInfos = append(participantInfos, models.ParticipantInfo{
								UserID:    participant.UniqueID,
								Username:  participant.Username,
								AvatarURL: participant.AvatarURL,
								IsOnline:  participant.IsActive,
							})
						} else {
							// If participant not found, still add them with basic info
							participantInfos = append(participantInfos, models.ParticipantInfo{
								UserID:    participantID,
								Username:  "Unknown User",
								AvatarURL: "",
								IsOnline:  false,
							})
						}
					}
					roomResponse.Participants = participantInfos

					roomList = append(roomList, roomResponse)
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{"rooms": roomList})
	}
}

// CreateRoomHandler creates a new study room
func CreateRoomHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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
		var req models.CreateRoomRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// First, try to find the user to get their unique_id
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		var user models.User
		userFilter := bson.M{"$or": []bson.M{
			{"_id": userIDStr},       // Try by MongoDB ObjectID
			{"unique_id": userIDStr}, // Try by unique_id
		}}
		err := users.FindOne(ctx, userFilter).Decode(&user)
		if err != nil {
			// If user not found, just use the original userIDStr
			user.UniqueID = userIDStr
		}

		now := time.Now()
		roomMap := bson.M{
			"name":             req.Name,
			"description":      req.Description,
			"type":             "shared",      // Set room type to shared
			"creator_id":       user.UniqueID, // Store user's unique_id
			"created_at":       now,
			"updated_at":       now,
			"last_activity_at": now,
			"is_active":        true,
			"max_participants": req.MaxParticipants,
			"participants":     []string{}, // Start with no participants - creator joins when they enter
			"materials":        []string{},
			"todos":            []string{},
			"notes":            []string{},
			"invitation_code":  "",
		}

		res, err := rooms.InsertOne(ctx, roomMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
			return
		}

		// Get the inserted ID
		var insertedID primitive.ObjectID
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			insertedID = oid
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get room ID"})
			return
		}

		// Construct the response directly instead of fetching from database
		roomResponse := models.RoomForResponse{
			ID:               insertedID.Hex(),
			Name:             req.Name,
			Description:      req.Description,
			Type:             "shared",
			CreatorID:        user.UniqueID, // Use user's unique_id
			CreatorUsername:  user.Username, // Include creator's username
			ParticipantCount: 0,             // Start with no participants - creator joins when they enter
			MaxParticipants:  req.MaxParticipants,
			MaterialsCount:   0,
			TodosCount:       0,
			NotesCount:       0,
			InvitationCode:   "",                         // Will be generated later if needed
			Participants:     []models.ParticipantInfo{}, // Start with no participants
			CreatedAt:        now,
			LastActivityAt:   now,
			IsActive:         true,
		}

		c.JSON(http.StatusCreated, gin.H{"room": roomResponse})
	}
}

// GetRoomHandler returns a specific room by ID
func GetRoomHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("id")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID required"})
			return
		}
		objID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}
		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var room models.Room
		err = rooms.FindOne(ctx, bson.M{"_id": objID}).Decode(&room)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Get creator's username
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		var creator models.User
		creatorFilter := bson.M{"unique_id": room.CreatorID}
		if err := users.FindOne(ctx, creatorFilter).Decode(&creator); err == nil {
			roomResponse := room.ToResponse()
			roomResponse.CreatorUsername = creator.Username

			// Populate participant details
			var participantInfos []models.ParticipantInfo
			for _, participantID := range room.Participants {
				var participant models.User
				participantFilter := bson.M{"unique_id": participantID}
				if err := users.FindOne(ctx, participantFilter).Decode(&participant); err == nil {
					participantInfos = append(participantInfos, models.ParticipantInfo{
						UserID:    participant.UniqueID,
						Username:  participant.Username,
						AvatarURL: participant.AvatarURL,
						IsOnline:  participant.IsActive,
					})
				} else {
					// If participant not found, still add them with basic info
					participantInfos = append(participantInfos, models.ParticipantInfo{
						UserID:    participantID,
						Username:  "Unknown User",
						AvatarURL: "",
						IsOnline:  false,
					})
				}
			}
			roomResponse.Participants = participantInfos

			c.JSON(http.StatusOK, gin.H{"room": roomResponse})
		} else {
			// If creator not found, still return room but with empty username
			roomResponse := room.ToResponse()

			// Populate participant details even if creator not found
			var participantInfos []models.ParticipantInfo
			for _, participantID := range room.Participants {
				var participant models.User
				participantFilter := bson.M{"unique_id": participantID}
				if err := users.FindOne(ctx, participantFilter).Decode(&participant); err == nil {
					participantInfos = append(participantInfos, models.ParticipantInfo{
						UserID:    participant.UniqueID,
						Username:  participant.Username,
						AvatarURL: participant.AvatarURL,
						IsOnline:  participant.IsActive,
					})
				} else {
					// If participant not found, still add them with basic info
					participantInfos = append(participantInfos, models.ParticipantInfo{
						UserID:    participantID,
						Username:  "Unknown User",
						AvatarURL: "",
						IsOnline:  false,
					})
				}
			}
			roomResponse.Participants = participantInfos

			c.JSON(http.StatusOK, gin.H{"room": roomResponse})
		}
	}
}

// UpdateRoomHandler updates a room
func UpdateRoomHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		roomID := c.Param("id")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID required"})
			return
		}
		roomObjID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}
		var req models.UpdateRoomRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		updateFields := bson.M{"updated_at": time.Now()}
		if req.Name != "" {
			updateFields["name"] = req.Name
		}
		if req.Description != "" {
			updateFields["description"] = req.Description
		}

		if req.MaxParticipants > 0 {
			updateFields["max_participants"] = req.MaxParticipants
		}

		update := bson.M{"$set": updateFields}
		res, err := rooms.UpdateOne(ctx, bson.M{"_id": roomObjID, "owner_id": userIDStr}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
			return
		}
		if res.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found or you don't have permission"})
			return
		}
		var room models.Room
		err = rooms.FindOne(ctx, bson.M{"_id": roomObjID}).Decode(&room)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated room"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"room": room.ToResponse()})
	}
}

// DeleteRoomHandler deletes a room
func DeleteRoomHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		roomID := c.Param("id")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID required"})
			return
		}
		roomObjID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}
		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Resolve the user's unique_id first
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		var user models.User
		userFilter := bson.M{"unique_id": userIDStr}
		_ = users.FindOne(ctx, userFilter).Decode(&user)
		if user.UniqueID == "" {
			user.UniqueID = userIDStr
		}

		// Only the creator (by unique_id) can delete
		res, err := rooms.DeleteOne(ctx, bson.M{"_id": roomObjID, "creator_id": user.UniqueID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
			return
		}
		if res.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found or you don't have permission"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
	}
}

// EnterRoomHandler allows users to enter a room they have access to (created or joined)
func EnterRoomHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		roomID := c.Param("id")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID required"})
			return
		}
		roomObjID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Find room by ID
		var room models.Room
		err = rooms.FindOne(ctx, bson.M{"_id": roomObjID, "is_active": true}).Decode(&room)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Check if user has access to this room (creator or already a participant)
		hasAccess := false
		if room.CreatorID == userIDStr {
			hasAccess = true
		} else {
			for _, participantID := range room.Participants {
				if participantID == userIDStr {
					hasAccess = true
					break
				}
			}
		}

		if !hasAccess {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this room"})
			return
		}

		// If user is not already a participant, add them
		isParticipant := false
		for _, participantID := range room.Participants {
			if participantID == userIDStr {
				isParticipant = true
				break
			}
		}

		if !isParticipant {
			// Check if room has space
			if len(room.Participants) >= room.MaxParticipants {
				c.JSON(http.StatusForbidden, gin.H{"error": "Room is full"})
				return
			}

			// Add user to participants
			update := bson.M{
				"$push": bson.M{"participants": userIDStr},
				"$set":  bson.M{"updated_at": time.Now()},
			}

			_, err = rooms.UpdateOne(ctx, bson.M{"_id": room.ID}, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enter room"})
				return
			}

			// Fetch updated room
			err = rooms.FindOne(ctx, bson.M{"_id": room.ID}).Decode(&room)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated room"})
				return
			}
		}

		// Update last activity
		_, err = rooms.UpdateOne(ctx, bson.M{"_id": room.ID}, bson.M{
			"$set": bson.M{"last_activity_at": time.Now()},
		})
		if err != nil {
			// Log error but don't fail the request
			fmt.Printf("Failed to update room activity: %v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{"room": room.ToResponse(), "message": "Successfully entered room"})
	}
}

// JoinRoomByCodeHandler allows users to join a room using an invitation code
func JoinRoomByCodeHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		var req models.JoinByCodeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Find room by invitation code
		var room models.Room
		err := rooms.FindOne(ctx, bson.M{"invitation_code": req.InvitationCode, "is_active": true}).Decode(&room)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid invitation code"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Check if user is already a participant
		for _, participantID := range room.Participants {
			if participantID == userIDStr {
				c.JSON(http.StatusConflict, gin.H{"error": "You are already a member of this room"})
				return
			}
		}

		// Check if room has space
		if len(room.Participants) >= room.MaxParticipants {
			c.JSON(http.StatusForbidden, gin.H{"error": "Room is full"})
			return
		}

		// Add user to participants
		update := bson.M{
			"$push": bson.M{"participants": userIDStr},
			"$set":  bson.M{"updated_at": time.Now()},
		}

		_, err = rooms.UpdateOne(ctx, bson.M{"_id": room.ID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join room"})
			return
		}

		// Fetch updated room
		err = rooms.FindOne(ctx, bson.M{"_id": room.ID}).Decode(&room)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated room"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"room": room.ToResponse(), "message": "Successfully joined room"})
	}
}

// GenerateInvitationCodeHandler generates a new invitation code for a room
func GenerateInvitationCodeHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		roomID := c.Param("id")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID required"})
			return
		}
		roomObjID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Verify user owns the room
		var room models.Room
		err = rooms.FindOne(ctx, bson.M{"_id": roomObjID, "creator_id": userIDStr}).Decode(&room)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found or you don't have permission"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Generate a unique invitation code
		invitationCode := generateUniqueInvitationCode()

		// Update room with new invitation code
		update := bson.M{
			"$set": bson.M{
				"invitation_code": invitationCode,
				"updated_at":      time.Now(),
			},
		}

		_, err = rooms.UpdateOne(ctx, bson.M{"_id": roomObjID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate invitation code"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"invitationCode": invitationCode})
	}
}

// LeaveRoomHandler allows a user to leave a room
func LeaveRoomHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		roomID := c.Param("id")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID required"})
			return
		}
		roomObjID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if user is a participant
		var room models.Room
		err = rooms.FindOne(ctx, bson.M{"_id": roomObjID}).Decode(&room)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Check if user is the creator (they can't leave, they should delete the room instead)
		if room.CreatorID == userIDStr {
			c.JSON(http.StatusForbidden, gin.H{"error": "Room creator cannot leave the room. Delete the room instead."})
			return
		}

		// Remove user from participants
		update := bson.M{
			"$pull": bson.M{"participants": userIDStr},
			"$set":  bson.M{"updated_at": time.Now()},
		}

		result, err := rooms.UpdateOne(ctx, bson.M{"_id": roomObjID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave room"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "You are not a member of this room"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully left the room"})
	}
}

// InviteUserToRoomHandler invites a user to join a room
func InviteUserToRoomHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		roomID := c.Param("id")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID required"})
			return
		}

		var req struct {
			TargetUserID string `json:"targetUserId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if room exists and user has permission to invite
		roomObjID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		var room models.Room
		if err := rooms.FindOne(ctx, bson.M{"_id": roomObjID}).Decode(&room); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		// Check if user is room owner or has permission to invite
		if room.CreatorID != userIDStr {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only room creators can invite users"})
			return
		}

		// Check if target user exists (by unique_id)
		var targetUser models.User
		if err := users.FindOne(ctx, bson.M{"unique_id": req.TargetUserID}).Decode(&targetUser); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Target user not found"})
			return
		}

		// Check if user is already in the room
		for _, participantID := range room.Participants {
			if participantID == targetUser.UniqueID {
				c.JSON(http.StatusConflict, gin.H{"error": "User is already in the room"})
				return
			}
		}

		// Get the inviter's username for the notification
		var inviter models.User
		if err := users.FindOne(ctx, bson.M{"unique_id": userIDStr}).Decode(&inviter); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inviter info"})
			return
		}

		// Create notification for the invited user (only for shared rooms interface)
		notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)
		roomInvitationNotification := models.CreateRoomInvitationNotification(
			targetUser.UniqueID,
			inviter.UniqueID,
			inviter.Username,
			roomID,
			room.Name,
		)

		_, err = notifications.InsertOne(ctx, roomInvitationNotification)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invitation notification"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User invited to room successfully"})
	}
}

// AcceptRoomInvitationHandler allows a user to accept a room invitation
func AcceptRoomInvitationHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		roomID := c.Param("id")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID required"})
			return
		}

		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if room exists
		roomObjID, err := primitive.ObjectIDFromHex(roomID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
			return
		}

		var room models.Room
		if err := rooms.FindOne(ctx, bson.M{"_id": roomObjID}).Decode(&room); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			return
		}

		// Check if user is already in the room
		for _, participantID := range room.Participants {
			if participantID == userIDStr {
				c.JSON(http.StatusConflict, gin.H{"error": "User is already in the room"})
				return
			}
		}

		// Check if room is at capacity
		if len(room.Participants) >= room.MaxParticipants {
			c.JSON(http.StatusConflict, gin.H{"error": "Room is at maximum capacity"})
			return
		}

		// Add user to room participants
		update := bson.M{
			"$push": bson.M{
				"participants": userIDStr,
			},
			"$set": bson.M{
				"updated_at": time.Now(),
			},
		}

		_, err = rooms.UpdateOne(ctx, bson.M{"_id": roomObjID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to room"})
			return
		}

		// Delete the invitation notification
		notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)
		_, err = notifications.DeleteMany(ctx, bson.M{
			"user_id":   userIDStr,
			"type":      models.NotificationTypeRoomInvitation,
			"target_id": roomID,
		})
		if err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: Failed to delete invitation notification: %v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully joined room"})
	}
}

// TestRoomHandler creates a test room to verify database connectivity
func TestRoomHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Try to insert a test document
		testRoom := bson.M{
			"name":             "TEST_ROOM",
			"description":      "Test room for debugging",
			"type":             "shared",
			"creator_id":       "test_user_123",
			"created_at":       time.Now(),
			"updated_at":       time.Now(),
			"last_activity_at": time.Now(),
			"is_active":        true,
			"max_participants": 4,
			"participants":     []string{"test_user_123"},
			"materials":        []string{},
			"todos":            []string{},
			"notes":            []string{},
			"invitation_code":  "",
		}

		res, err := rooms.InsertOne(ctx, testRoom)
		if err != nil {
			fmt.Printf("TestRoomHandler: insert error=%v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert test room", "details": err.Error()})
			return
		}

		// Try to find the test room
		var foundRoom bson.M
		err = rooms.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&foundRoom)
		if err != nil {
			fmt.Printf("TestRoomHandler: find error=%v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find test room", "details": err.Error()})
			return
		}

		// Clean up - delete the test room
		_, err = rooms.DeleteOne(ctx, bson.M{"_id": res.InsertedID})
		if err != nil {
			fmt.Printf("TestRoomHandler: delete error=%v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message":     "Database test successful",
			"inserted_id": res.InsertedID,
			"found_room":  foundRoom,
		})
	}
}

// DebugListAllRoomsHandler lists all rooms in the database (for debugging only)
func DebugListAllRoomsHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get all rooms without any filter
		cursor, err := rooms.Find(ctx, bson.M{})
		if err != nil {
			fmt.Printf("DebugListAllRoomsHandler: database error=%v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
			return
		}
		defer cursor.Close(ctx)

		var allRooms []bson.M
		for cursor.Next(ctx) {
			var room bson.M
			if err := cursor.Decode(&room); err == nil {
				allRooms = append(allRooms, room)
			} else {
				fmt.Printf("DebugListAllRoomsHandler: error decoding room=%v\n", err)
			}
		}

		fmt.Printf("DebugListAllRoomsHandler: found %d total rooms in database\n", len(allRooms))

		c.JSON(http.StatusOK, gin.H{
			"total_rooms": len(allRooms),
			"rooms":       allRooms,
		})
	}
}

// SimpleTestHandler tests basic MongoDB operations without authentication
func SimpleTestHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		rooms := mongoClient.GetCollection(database.CollectionNames.Rooms)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Test 1: Count total rooms
		totalRooms, err := rooms.CountDocuments(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to count rooms",
				"details": err.Error(),
			})
			return
		}

		// Test 2: Try to find any room
		var sampleRoom bson.M
		err = rooms.FindOne(ctx, bson.M{}).Decode(&sampleRoom)
		hasRooms := err == nil

		// Test 3: Try to insert a test document
		testDoc := bson.M{
			"test":      true,
			"timestamp": time.Now(),
		}

		res, err := rooms.InsertOne(ctx, testDoc)
		insertSuccess := err == nil
		var insertedID interface{}
		if insertSuccess {
			insertedID = res.InsertedID
			// Clean up - delete the test document
			rooms.DeleteOne(ctx, bson.M{"_id": res.InsertedID})
		}

		c.JSON(http.StatusOK, gin.H{
			"message":             "MongoDB connectivity test",
			"total_rooms":         totalRooms,
			"has_rooms":           hasRooms,
			"insert_test_success": insertSuccess,
			"inserted_id":         insertedID,
			"sample_room":         sampleRoom,
		})
	}
}

// Helper function to generate unique invitation codes
func generateUniqueInvitationCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const codeLength = 8

	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(code)
}
