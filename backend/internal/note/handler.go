package note

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/studyplatform/backend/pkg/database"
	"github.com/studyplatform/backend/pkg/models"
)

// ListNotesHandler returns all notes for the authenticated user
func ListNotesHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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
		// userID is now UniqueID (e.g., "USER123456"), not MongoDB ObjectID
		notes := mongoClient.GetCollection(database.CollectionNames.Notes)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		filter := bson.M{"creator_id": userIDStr} // Use creator_id instead of owner_id
		cursor, err := notes.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer cursor.Close(ctx)
		var noteList []models.NoteForResponse
		for cursor.Next(ctx) {
			var note models.Note
			if err := cursor.Decode(&note); err == nil {
				noteList = append(noteList, note.ToResponse())
			}
		}
		c.JSON(http.StatusOK, gin.H{"notes": noteList})
	}
}

// CreateNoteHandler creates a new note
func CreateNoteHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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
		// Don't convert to ObjectID - use UniqueID directly

		var req models.CreateNoteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		notes := mongoClient.GetCollection(database.CollectionNames.Notes)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		note := models.Note{
			Title:        req.Content[:int(math.Min(float64(len(req.Content)), 50))], // Use first 50 chars as title
			Content:      req.Content,
			RoomID:       req.RoomID,
			CreatorID:    userIDStr, // Use UniqueID directly
			Tags:         []string{},
			SharedWith:   []string{},
			IsPublic:     req.IsShared, // Map IsShared to IsPublic
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			LastEditedBy: userIDStr,
		}

		res, err := notes.InsertOne(ctx, note)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
			return
		}

		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			note.ID = oid
		}

		c.JSON(http.StatusCreated, gin.H{"note": note.ToResponse()})
	}
}

// GetNotesHandler retrieves all notes for the authenticated user
func GetNotesHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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
		// Don't convert to ObjectID - use UniqueID directly

		notes := mongoClient.GetCollection(database.CollectionNames.Notes)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get notes created by the user
		filter := bson.M{"creator_id": userIDStr} // Use UniqueID directly
		cursor, err := notes.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
			return
		}
		defer cursor.Close(ctx)

		var notesList []models.Note
		if err = cursor.All(ctx, &notesList); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode notes"})
			return
		}

		// Convert to response format
		var responseNotes []models.NoteForResponse
		for _, note := range notesList {
			responseNotes = append(responseNotes, note.ToResponse())
		}

		c.JSON(http.StatusOK, gin.H{"notes": responseNotes})
	}
}

// GetRoomNotesHandler retrieves notes for a specific room
func GetRoomNotesHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("id")
		fmt.Printf("GetRoomNotesHandler: roomID = '%s', all params = %+v\n", roomID, c.Params)
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID required"})
			return
		}

		notes := mongoClient.GetCollection(database.CollectionNames.Notes)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get notes for the specific room
		filter := bson.M{"room_id": roomID}
		cursor, err := notes.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch room notes"})
			return
		}
		defer cursor.Close(ctx)

		var notesList []models.Note
		if err = cursor.All(ctx, &notesList); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode room notes"})
			return
		}

		// Convert to response format
		var responseNotes []models.NoteForResponse
		for _, note := range notesList {
			responseNotes = append(responseNotes, note.ToResponse())
		}

		c.JSON(http.StatusOK, gin.H{"notes": responseNotes})
	}
}

// UpdateNoteHandler updates an existing note
func UpdateNoteHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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
		// Don't convert to ObjectID - use UniqueID directly

		noteID := c.Param("id")
		if noteID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Note ID required"})
			return
		}
		noteObjID, err := primitive.ObjectIDFromHex(noteID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
			return
		}
		var req models.UpdateNoteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		notes := mongoClient.GetCollection(database.CollectionNames.Notes)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		updateFields := bson.M{"updated_at": time.Now()}
		if req.Content != "" {
			updateFields["content"] = req.Content
		}
		if req.IsShared != nil {
			updateFields["is_public"] = *req.IsShared // Map IsShared to is_public
		}
		update := bson.M{"$set": updateFields}
		res, err := notes.UpdateOne(ctx, bson.M{"_id": noteObjID, "creator_id": userIDStr}, update) // Use UniqueID directly
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
			return
		}
		if res.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or you don't have permission"})
			return
		}
		var note models.Note
		err = notes.FindOne(ctx, bson.M{"_id": noteObjID}).Decode(&note)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated note"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"note": note.ToResponse()})
	}
}

// DeleteNoteHandler deletes a note
func DeleteNoteHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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
		// Don't convert to ObjectID - use UniqueID directly

		noteID := c.Param("id")
		if noteID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Note ID required"})
			return
		}
		noteObjID, err := primitive.ObjectIDFromHex(noteID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
			return
		}
		notes := mongoClient.GetCollection(database.CollectionNames.Notes)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := notes.DeleteOne(ctx, bson.M{"_id": noteObjID, "creator_id": userIDStr}) // Use UniqueID directly
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
			return
		}
		if res.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found or you don't have permission"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
	}
}
