package notification

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/studyplatform/backend/pkg/database"
	"github.com/studyplatform/backend/pkg/models"
)

// ListNotificationsHandler returns user's notifications
func ListNotificationsHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := notifications.Find(ctx, bson.M{"user_id": userIDStr})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer cursor.Close(ctx)

		var notificationList []models.NotificationForResponse
		for cursor.Next(ctx) {
			var notification models.Notification
			if err := cursor.Decode(&notification); err == nil {
				notificationList = append(notificationList, notification.ToResponse())
			}
		}

		c.JSON(http.StatusOK, gin.H{"notifications": notificationList})
	}
}

// MarkNotificationReadHandler marks a notification as read and optionally deletes it
func MarkNotificationReadHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		notificationID := c.Param("id")
		if notificationID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID required"})
			return
		}

		objID, err := primitive.ObjectIDFromHex(notificationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
			return
		}

		notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Delete the notification instead of just marking it as read
		res, err := notifications.DeleteOne(ctx, bson.M{"_id": objID, "user_id": userIDStr})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
			return
		}

		if res.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
	}
}

// CreateNotificationHandler creates a new notification (internal use)
func CreateNotificationHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		_, ok := userID.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID in context"})
			return
		}

		var req struct {
			UserID   string `json:"userId" binding:"required"`
			Type     string `json:"type" binding:"required"`
			Title    string `json:"title" binding:"required"`
			Message  string `json:"message" binding:"required"`
			TargetID string `json:"targetId"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		notificationMap := bson.M{
			"user_id":    req.UserID,
			"type":       req.Type,
			"title":      req.Title,
			"message":    req.Message,
			"target_id":  req.TargetID,
			"is_read":    false,
			"created_at": time.Now(),
		}

		res, err := notifications.InsertOne(ctx, notificationMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
			return
		}

		var notification models.Notification
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			notification.ID = oid
		}
		err = notifications.FindOne(ctx, bson.M{"_id": notification.ID}).Decode(&notification)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created notification"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"notification": notification.ToResponse()})
	}
}

// ClearFriendRequestNotificationsHandler clears all friend request notifications for a user (for debugging)
func ClearFriendRequestNotificationsHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Delete all friend request notifications for the user
		res, err := notifications.DeleteMany(ctx, bson.M{
			"user_id": userIDStr,
			"type":    "friend_request",
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear notifications"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "Friend request notifications cleared",
			"deleted_count": res.DeletedCount,
		})
	}
}

// ClearAllNotificationsHandler clears ALL notifications for a user (NUCLEAR OPTION)
func ClearAllNotificationsHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Delete ALL notifications for the user
		res, err := notifications.DeleteMany(ctx, bson.M{
			"user_id": userIDStr,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear all notifications"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "ALL notifications cleared",
			"deleted_count": res.DeletedCount,
		})
	}
}

// DeleteNotificationHandler deletes a specific notification
func DeleteNotificationHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		notificationID := c.Param("id")
		if notificationID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Notification ID required"})
			return
		}

		objID, err := primitive.ObjectIDFromHex(notificationID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
			return
		}

		notifications := mongoClient.GetCollection(database.CollectionNames.Notifications)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := notifications.DeleteOne(ctx, bson.M{"_id": objID, "user_id": userIDStr})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
			return
		}

		if res.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
	}
}
