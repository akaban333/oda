package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/studyplatform/backend/pkg/database"
	"github.com/studyplatform/backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SendFriendRequestHandler handles sending a friend request by unique ID
func SendFriendRequestHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var req struct {
			UniqueID string `json:"uniqueId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		if req.UniqueID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unique ID is required"})
			return
		}

		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Find the requesting user by UniqueID (from JWT token)
		var me models.User
		if err := users.FindOne(ctx, bson.M{"unique_id": userID.(string)}).Decode(&me); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		// Find the target user by unique ID
		var target models.User
		if err := users.FindOne(ctx, bson.M{"unique_id": req.UniqueID}).Decode(&target); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Target user not found"})
			return
		}

		if target.ID == me.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add yourself as a friend"})
			return
		}

		// Check if already friends or pending
		for _, f := range me.Friends {
			if f.UserID == target.UniqueID && (f.Status == "accepted" || f.Status == "pending") {
				c.JSON(http.StatusConflict, gin.H{"error": "Already friends or request pending"})
				return
			}
		}

		// Initialize Friends array if it's nil
		if me.Friends == nil {
			me.Friends = []models.Friend{}
		}
		if target.Friends == nil {
			target.Friends = []models.Friend{}
		}

		// Add friend request to requesting user (using UniqueID)
		_, err := users.UpdateOne(ctx, bson.M{"_id": me.ID}, bson.M{
			"$push": bson.M{"friends": models.Friend{UserID: target.UniqueID, Status: "pending", Since: time.Now()}},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
			return
		}

		// Add friend request to target user (using UniqueID)
		_, err = users.UpdateOne(ctx, bson.M{"_id": target.ID}, bson.M{
			"$push": bson.M{"friends": models.Friend{UserID: me.UniqueID, Status: "requested", Since: time.Now()}},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update target user"})
			return
		}

		// Create notification for the target user
		notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)
		friendRequestNotification := models.CreateFriendRequestNotification(target.ID.Hex(), me.ID.Hex(), me.Username)

		_, err = notifications.InsertOne(ctx, friendRequestNotification)
		if err != nil {
			// Don't fail the entire request if notification creation fails
		}

		c.JSON(http.StatusOK, gin.H{"message": "Friend request sent"})
	}
}

// AcceptFriendRequestHandler handles accepting a friend request
func AcceptFriendRequestHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		friendID := c.Param("id")
		if friendID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Friend ID required"})
			return
		}
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Find current user by UniqueID
		var me models.User
		if err := users.FindOne(ctx, bson.M{"unique_id": userID.(string)}).Decode(&me); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		// Find friend by UniqueID
		var friend models.User
		if err := users.FindOne(ctx, bson.M{"unique_id": friendID}).Decode(&friend); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
			return
		}

		// Update both users' friend status to accepted (using UniqueIDs)
		_, err := users.UpdateOne(ctx, bson.M{"_id": me.ID, "friends.user_id": friendID}, bson.M{"$set": bson.M{"friends.$.status": "accepted", "friends.$.since": time.Now()}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept request"})
			return
		}
		_, err = users.UpdateOne(ctx, bson.M{"_id": friend.ID, "friends.user_id": me.UniqueID}, bson.M{"$set": bson.M{"friends.$.status": "accepted", "friends.$.since": time.Now()}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update friend"})
			return
		}

		// Create notification for the friend that their request was accepted
		notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)
		friendAcceptedNotification := models.CreateFriendAcceptedNotification(friend.ID.Hex(), me.ID.Hex(), me.Username)
		_, err = notifications.InsertOne(ctx, friendAcceptedNotification)
		if err != nil {
			// Don't fail if notification creation fails
		}

		// Remove friend request notifications between these two users
		_, err = notifications.DeleteMany(ctx, bson.M{
			"type": "friend_request",
			"$or": []bson.M{
				{"user_id": me.ID.Hex(), "data.requester_id": friend.ID.Hex()},
				{"user_id": friend.ID.Hex(), "data.requester_id": me.ID.Hex()},
			},
		})
		if err != nil {
			// Don't fail if notification cleanup fails
		}

		c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted"})
	}
}

// RejectFriendRequestHandler handles rejecting a friend request
func RejectFriendRequestHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		friendID := c.Param("id")
		if friendID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Friend ID required"})
			return
		}
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Find current user by UniqueID
		var me models.User
		if err := users.FindOne(ctx, bson.M{"unique_id": userID.(string)}).Decode(&me); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		// Find friend by UniqueID
		var friend models.User
		if err := users.FindOne(ctx, bson.M{"unique_id": friendID}).Decode(&friend); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
			return
		}

		// Remove friend entries from both users (using UniqueIDs)
		_, err := users.UpdateOne(ctx, bson.M{"_id": me.ID}, bson.M{"$pull": bson.M{"friends": bson.M{"user_id": friendID}}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject request"})
			return
		}
		_, err = users.UpdateOne(ctx, bson.M{"_id": friend.ID}, bson.M{"$pull": bson.M{"friends": bson.M{"user_id": me.UniqueID}}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update friend"})
			return
		}

		// Remove friend request notifications between these two users
		notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)
		_, err = notifications.DeleteMany(ctx, bson.M{
			"type": "friend_request",
			"$or": []bson.M{
				{"user_id": me.ID.Hex(), "data.requester_id": friend.ID.Hex()},
				{"user_id": friend.ID.Hex(), "data.requester_id": me.ID.Hex()},
			},
		})
		if err != nil {
			// Don't fail if notification cleanup fails
		}

		c.JSON(http.StatusOK, gin.H{"message": "Friend request rejected"})
	}
}

// ListFriendsHandler returns the list of all friend relationships
func ListFriendsHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userIDStr := userID.(string)

		var me models.User
		userFilter := bson.M{"unique_id": userIDStr}
		if err := users.FindOne(ctx, userFilter).Decode(&me); err != nil {
			// Try alternative lookup methods
			var altUser models.User
			if err := users.FindOne(ctx, bson.M{"_id": userIDStr}).Decode(&altUser); err == nil {
				me = altUser
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found", "details": "Could not find user with unique_id or _id: " + userIDStr})
				return
			}
		}

		var friendsList []map[string]interface{}
		for _, f := range me.Friends {
			var friend models.User
			var found bool

			// Try to find friend by UniqueID first (new format)
			if err := users.FindOne(ctx, bson.M{"unique_id": f.UserID}).Decode(&friend); err == nil {
				found = true
			} else {
				// Fallback: try to find by ObjectID (old format)
				if objID, err := primitive.ObjectIDFromHex(f.UserID); err == nil {
					if err := users.FindOne(ctx, bson.M{"_id": objID}).Decode(&friend); err == nil {
						found = true
					}
				}
			}

			if found {
				friendData := map[string]interface{}{
					"id":          f.UserID,
					"userId":      f.UserID,
					"status":      f.Status,
					"since":       f.Since,
					"username":    friend.Username,
					"uniqueId":    friend.UniqueID,
					"firstName":   friend.FirstName,
					"lastName":    friend.LastName,
					"displayName": friend.FirstName + " " + friend.LastName,
					"avatarUrl":   friend.AvatarURL,
				}
				friendsList = append(friendsList, friendData)
			}
		}

		c.JSON(http.StatusOK, gin.H{"friends": friendsList})
	}
}

// ListFriendRequestsHandler returns the list of pending/received friend requests
func ListFriendRequestsHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		userIDStr := userID.(string)
		var me models.User
		userFilter := bson.M{"unique_id": userIDStr}
		if err := users.FindOne(ctx, userFilter).Decode(&me); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}
		var sent, received []models.UserForResponse
		for _, f := range me.Friends {
			var friend models.User
			var found bool

			// Try to find friend by UniqueID first (new format)
			if err := users.FindOne(ctx, bson.M{"unique_id": f.UserID}).Decode(&friend); err == nil {
				found = true
			} else {
				// Fallback: try to find by ObjectID (old format)
				if objID, err := primitive.ObjectIDFromHex(f.UserID); err == nil {
					if err := users.FindOne(ctx, bson.M{"_id": objID}).Decode(&friend); err == nil {
						found = true
					}
				}
			}

			if found {
				if f.Status == "pending" {
					sent = append(sent, friend.ToResponse())
				} else if f.Status == "requested" {
					received = append(received, friend.ToResponse())
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{"sent": sent, "received": received})
	}
}

// SearchUsersHandler allows searching users by unique ID only
func SearchUsersHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("search")

		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
			return
		}

		// Trim whitespace and validate the query
		query = strings.TrimSpace(query)

		if len(query) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Search query must be at least 3 characters long"})
			return
		}

		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Only search by unique_id (exact match for user ID)
		filter := bson.M{"unique_id": query}

		cursor, err := users.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer cursor.Close(ctx)

		var results []models.UserForResponse
		for cursor.Next(ctx) {
			var user models.User
			if err := cursor.Decode(&user); err == nil {
				userResponse := user.ToResponse()
				results = append(results, userResponse)
			}
		}

		response := gin.H{"results": results}
		c.JSON(http.StatusOK, response)
	}
}

// GetUserByIDHandler returns a user profile by unique ID
func GetUserByIDHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID required"})
			return
		}
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var user models.User
		// Try by ObjectID first, then by unique_id
		objID, err := primitive.ObjectIDFromHex(id)
		if err == nil {
			if err := users.FindOne(ctx, bson.M{"_id": objID}).Decode(&user); err == nil {
				c.JSON(http.StatusOK, gin.H{"user": user.ToResponse()})
				return
			}
		}
		if err := users.FindOne(ctx, bson.M{"unique_id": id}).Decode(&user); err == nil {
			c.JSON(http.StatusOK, gin.H{"user": user.ToResponse()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}
}

// RemoveFriendHandler handles removing a friend
func RemoveFriendHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		friendID := c.Param("id")
		if friendID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Friend ID required"})
			return
		}
		users := mongoClient.GetCollection(database.CollectionNames.Users)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Find current user by UniqueID
		var me models.User
		if err := users.FindOne(ctx, bson.M{"unique_id": userID.(string)}).Decode(&me); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		// Find friend by UniqueID
		var friend models.User
		if err := users.FindOne(ctx, bson.M{"unique_id": friendID}).Decode(&friend); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
			return
		}

		// Remove friend entries from both users (using UniqueIDs)
		_, err := users.UpdateOne(ctx, bson.M{"_id": me.ID}, bson.M{"$pull": bson.M{"friends": bson.M{"user_id": friendID}}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove friend"})
			return
		}
		_, err = users.UpdateOne(ctx, bson.M{"_id": friend.ID}, bson.M{"$pull": bson.M{"friends": bson.M{"user_id": me.UniqueID}}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update friend"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Friend removed successfully"})
	}
}

// MigrateFriendsToUniqueID converts existing ObjectID-based friend relationships to UniqueID-based ones
func MigrateFriendsToUniqueID(mongoClient *database.MongoClient) error {
	users := mongoClient.GetCollection(database.CollectionNames.Users)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Find all users with friends
	cursor, err := users.Find(ctx, bson.M{"friends": bson.M{"$exists": true, "$ne": []models.Friend{}}})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			continue
		}

		var updated bool
		for i, friend := range user.Friends {
			// Check if this is an ObjectID (24 character hex string)
			if len(friend.UserID) == 24 {
				// Try to convert to ObjectID and find the user
				if objID, err := primitive.ObjectIDFromHex(friend.UserID); err == nil {
					var friendUser models.User
					if err := users.FindOne(ctx, bson.M{"_id": objID}).Decode(&friendUser); err == nil {
						// Update the friend relationship to use UniqueID
						user.Friends[i].UserID = friendUser.UniqueID
						updated = true
					}
				}
			}
		}

		// Update the user if any friends were migrated
		if updated {
			_, err := users.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"friends": user.Friends}})
			if err != nil {
				// Log error but continue with other users
				fmt.Printf("Failed to migrate friends for user %s: %v\n", user.Username, err)
			}
		}
	}

	return nil
}
