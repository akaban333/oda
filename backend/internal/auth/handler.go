package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/uuid"
	"github.com/studyplatform/backend/pkg/auth"
	"github.com/studyplatform/backend/pkg/database"
	"github.com/studyplatform/backend/pkg/logger"
	"github.com/studyplatform/backend/pkg/models"
)

// RegisterHandler handles user registration
func RegisterHandler(mongoClient *database.MongoClient, jwtManager *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SignupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if email already exists (username can be duplicated)
		filter := bson.M{"email": strings.ToLower(req.Email)}
		count, err := users.CountDocuments(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}

		// Hash password
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		uniqueID := uuid.NewString()
		userMap := bson.M{
			"unique_id":   uniqueID,
			"username":    req.Username,
			"email":       strings.ToLower(req.Email),
			"password":    hashedPassword,
			"first_name":  "", // Default empty string
			"last_name":   "", // Default empty string
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
			"is_active":   true,
			"is_verified": false,
		}
		res, err := users.InsertOne(ctx, userMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		var user models.User
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			user.ID = oid
		}
		// Fetch the user from DB to get all fields
		err = users.FindOne(ctx, bson.M{"_id": user.ID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created user"})
			return
		}
		accessToken, err := jwtManager.GenerateAccessToken(user.UniqueID, user.Username, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}
		refreshToken, err := jwtManager.GenerateRefreshToken(user.UniqueID, user.Username, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"user":         user.ToResponse(),
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}

// LoginHandler handles user login
func LoginHandler(mongoClient *database.MongoClient, jwtManager *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.User
		err := users.FindOne(ctx, bson.M{"email": strings.ToLower(req.Email)}).Decode(&user)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Debug: Check what password field exists
		var rawUser bson.M
		err = users.FindOne(ctx, bson.M{"email": strings.ToLower(req.Email)}).Decode(&rawUser)
		if err == nil {
			logger.Info("User found", logger.Field("email", req.Email), logger.Field("hasPassword", rawUser["password"] != nil), logger.Field("hasPasswordHash", rawUser["password_hash"] != nil))
		}

		// Try both password field names
		var passwordToCheck string
		if user.PasswordHash != "" {
			passwordToCheck = user.PasswordHash
		} else if rawUser["password"] != nil {
			passwordToCheck = rawUser["password"].(string)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		if err := auth.ComparePassword(passwordToCheck, req.Password); err != nil {
			logger.Warn("Password comparison failed", logger.Field("error", err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		accessToken, err := jwtManager.GenerateAccessToken(user.UniqueID, user.Username, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}
		refreshToken, err := jwtManager.GenerateRefreshToken(user.UniqueID, user.Username, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":         user.ToResponse(),
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}

// MeHandler returns the current authenticated user's profile
func MeHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var user models.User
		// Try to find user by either MongoDB ObjectID or unique_id
		userFilter := bson.M{"$or": []bson.M{
			{"_id": userIDStr},       // Try by MongoDB ObjectID
			{"unique_id": userIDStr}, // Try by unique_id
		}}
		err := users.FindOne(ctx, userFilter).Decode(&user)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user.ToResponse()})
	}
}

// UpdateProfileHandler allows the authenticated user to update their profile (bio, avatar, etc.)
func UpdateProfileHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		// Parse raw body first
		raw, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}
		var rawMap map[string]interface{}
		if err := json.Unmarshal(raw, &rawMap); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Debug logging
		fmt.Printf("DEBUG: UpdateProfile - Raw request: %+v\n", rawMap)
		fmt.Printf("DEBUG: UpdateProfile - User ID: %s\n", userIDStr)

		updateFields := bson.M{"updated_at": time.Now()}

		// Handle username update - no uniqueness check needed
		if v, ok := rawMap["username"]; ok {
			if str, ok := v.(string); ok {
				updateFields["username"] = str
			}
		}

		// Handle firstName - allow empty strings
		if v, ok := rawMap["firstName"]; ok {
			if str, ok := v.(string); ok {
				updateFields["first_name"] = str
			}
		}

		// Handle lastName - allow empty strings
		if v, ok := rawMap["lastName"]; ok {
			if str, ok := v.(string); ok {
				updateFields["last_name"] = str
			}
		}

		// Handle bio - allow empty strings
		if v, ok := rawMap["bio"]; ok {
			if str, ok := v.(string); ok {
				updateFields["bio"] = str
			}
		}

		// Handle avatarUrl - allow empty strings
		if v, ok := rawMap["avatarUrl"]; ok {
			if str, ok := v.(string); ok {
				updateFields["avatar_url"] = str
			}
		}
		if len(updateFields) == 1 { // only updated_at
			c.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
			return
		}
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		update := bson.M{"$set": updateFields}
		fmt.Printf("DEBUG: UpdateProfile - Final update fields: %+v\n", updateFields)
		// Try to find user by either MongoDB ObjectID or unique_id
		userFilter := bson.M{"$or": []bson.M{
			{"_id": userIDStr},       // Try by MongoDB ObjectID
			{"unique_id": userIDStr}, // Try by unique_id
		}}
		res, err := users.UpdateOne(ctx, userFilter, update)
		fmt.Printf("DEBUG: Update profile for user %s, matched: %d, modified: %d, err: %v\n", userIDStr, res.MatchedCount, res.ModifiedCount, err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}
		var user models.User
		err = users.FindOne(ctx, bson.M{"unique_id": userIDStr}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user.ToResponse()})
	}
}

// MigrateAddUniqueIDToUsers assigns a unique_id to all users missing it
func MigrateAddUniqueIDToUsers(mongoClient *database.MongoClient) error {
	users := mongoClient.GetCollection(database.CollectionNames.Users)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cursor, err := users.Find(ctx, bson.M{"unique_id": bson.M{"$exists": false}})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	count := 0
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		uniqueID := uuid.NewString()
		_, err := users.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"unique_id": uniqueID}})
		if err == nil {
			count++
		}
	}
	println("Migration complete. Users updated:", count)
	return nil
}

// MigratePasswordFields updates existing users to use password_hash field
func MigratePasswordFields(mongoClient *database.MongoClient) error {
	users := mongoClient.GetCollection(database.CollectionNames.Users)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find users with old password field
	cursor, err := users.Find(ctx, bson.M{"password": bson.M{"$exists": true}})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var updatedCount int
	for cursor.Next(ctx) {
		var user bson.M
		if err := cursor.Decode(&user); err != nil {
			continue
		}

		// Update user to use password_hash field
		update := bson.M{
			"$set": bson.M{
				"password_hash": user["password"],
				"updated_at":    time.Now(),
			},
			"$unset": bson.M{"password": ""},
		}

		_, err := users.UpdateOne(ctx, bson.M{"_id": user["_id"]}, update)
		if err != nil {
			logger.Warn("Failed to migrate user password", logger.Field("user_id", user["_id"]), logger.Field("error", err))
		} else {
			updatedCount++
		}
	}

	logger.Info("Password migration complete", logger.Field("users_updated", updatedCount))
	return nil
}

// MigrateRemoveUsernameUniqueness removes any unique constraints on username field
func MigrateRemoveUsernameUniqueness(mongoClient *database.MongoClient) error {
	users := mongoClient.GetCollection(database.CollectionNames.Users)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Drop any existing indexes on username field
	indexes, err := users.Indexes().List(ctx)
	if err != nil {
		logger.Warn("Failed to list indexes", logger.Field("error", err))
		return err
	}

	// Check if there's a unique index on username
	for indexes.Next(ctx) {
		var index bson.M
		if err := indexes.Decode(&index); err != nil {
			continue
		}

		// Check if this is a unique index on username
		if key, ok := index["key"].(bson.M); ok {
			if _, hasUsername := key["username"]; hasUsername {
				if unique, ok := index["unique"].(bool); ok && unique {
					logger.Info("Found unique index on username, dropping it", logger.Field("index_name", index["name"]))
					// Drop the unique index
					if _, err := users.Indexes().DropOne(ctx, index["name"].(string)); err != nil {
						logger.Warn("Failed to drop unique index on username", logger.Field("error", err))
					}
				}
			}
		}
	}

	logger.Info("Username uniqueness migration completed")
	return nil
}

// LogoutHandler handles user logout
func LogoutHandler(mongoClient *database.MongoClient, jwtManager *auth.Manager) gin.HandlerFunc {
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
		// Invalidate refresh tokens for the user
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var err error
		_, err = users.UpdateOne(ctx, bson.M{"unique_id": userIDStr}, bson.M{"$set": bson.M{"refresh_tokens": []models.RefreshToken{}}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

// RefreshTokenHandler handles token refresh
func RefreshTokenHandler(mongoClient *database.MongoClient, jwtManager *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			RefreshToken string `json:"refreshToken" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		claims, err := jwtManager.ValidateRefreshToken(req.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}
		// Generate new access token
		accessToken, err := jwtManager.GenerateAccessToken(claims.UserID, claims.Username, claims.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
	}
}

// UpdateXPHandler updates user's XP
func UpdateXPHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		var req models.UpdateXPRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Update user's XP
		update := bson.M{
			"$inc": bson.M{"xp": req.XP},
			"$set": bson.M{"updated_at": time.Now()},
		}

		result, err := users.UpdateOne(ctx, bson.M{"unique_id": userIDStr}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update XP"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Record XP history
		xpHistory := models.XPHistory{
			ID:        primitive.NewObjectID(),
			UserID:    userIDStr,
			XP:        req.XP,
			Source:    req.Source,
			SessionID: req.SessionID,
			CreatedAt: time.Now(),
		}

		xpHistoryCollection := mongoClient.GetCollection("xp_history")
		_, err = xpHistoryCollection.InsertOne(ctx, xpHistory)
		if err != nil {
			logger.Warn("Failed to record XP history", logger.Field("error", err))
			// Don't fail the request if XP history recording fails
		}

		c.JSON(http.StatusOK, gin.H{"message": "XP updated successfully", "xp_earned": req.XP})
	}
}
