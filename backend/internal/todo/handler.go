package todo

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

// ListTodosHandler returns all todos for the authenticated user
func ListTodosHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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
		todos := mongoClient.GetCollection(database.CollectionNames.Todos)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		filter := bson.M{"$or": []bson.M{
			{"creator_id": userIDStr},
			{"assignee_ids": userIDStr},
		}}
		cursor, err := todos.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer cursor.Close(ctx)
		var todoList []models.TodoForResponse
		for cursor.Next(ctx) {
			var todo models.Todo
			if err := cursor.Decode(&todo); err == nil {
				todoList = append(todoList, todo.ToResponse())
			}
		}
		c.JSON(http.StatusOK, gin.H{"todos": todoList})
	}
}

// CreateTodoHandler creates a new todo item
func CreateTodoHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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
		var req models.CreateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		todos := mongoClient.GetCollection(database.CollectionNames.Todos)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		todoMap := bson.M{
			"title":        req.Title,
			"description":  req.Description,
			"completed":    false,
			"due_date":     req.DueDate,
			"priority":     req.Priority,
			"room_id":      req.RoomID,
			"creator_id":   userIDStr,
			"assignee_ids": req.AssigneeIDs,
			"tags":         req.Tags,
			"created_at":   time.Now(),
			"updated_at":   time.Now(),
		}
		res, err := todos.InsertOne(ctx, todoMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
			return
		}

		// Get the inserted ID
		insertedID, ok := res.InsertedID.(primitive.ObjectID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inserted ID"})
			return
		}

		// Fetch the created todo to return it
		var todo models.Todo
		err = todos.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created todo"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"todo": todo.ToResponse()})
	}
}

// GetTodoHandler returns a specific todo by ID
func GetTodoHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		todoID := c.Param("id")
		if todoID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Todo ID required"})
			return
		}
		objID, err := primitive.ObjectIDFromHex(todoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
			return
		}
		todos := mongoClient.GetCollection(database.CollectionNames.Todos)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var todo models.Todo
		err = todos.FindOne(ctx, bson.M{"_id": objID}).Decode(&todo)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"todo": todo.ToResponse()})
	}
}

// UpdateTodoHandler updates a todo item
func UpdateTodoHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		todoID := c.Param("id")
		if todoID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Todo ID required"})
			return
		}
		todoObjID, err := primitive.ObjectIDFromHex(todoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
			return
		}
		var req models.UpdateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		todos := mongoClient.GetCollection(database.CollectionNames.Todos)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		updateFields := bson.M{"updated_at": time.Now()}
		if req.Title != "" {
			updateFields["title"] = req.Title
		}
		if req.Description != "" {
			updateFields["description"] = req.Description
		}
		if req.DueDate != nil {
			updateFields["due_date"] = req.DueDate
		}
		if req.Priority > 0 {
			updateFields["priority"] = req.Priority
		}
		if len(req.AssigneeIDs) > 0 {
			updateFields["assignee_ids"] = req.AssigneeIDs
		}
		if len(req.Tags) > 0 {
			updateFields["tags"] = req.Tags
		}
		update := bson.M{"$set": updateFields}
		res, err := todos.UpdateOne(ctx, bson.M{"_id": todoObjID, "creator_id": userIDStr}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
			return
		}
		if res.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or you don't have permission"})
			return
		}
		var todo models.Todo
		err = todos.FindOne(ctx, bson.M{"_id": todoObjID}).Decode(&todo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated todo"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"todo": todo.ToResponse()})
	}
}

// DeleteTodoHandler deletes a todo item
func DeleteTodoHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		todoID := c.Param("id")
		if todoID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Todo ID required"})
			return
		}
		objID, err := primitive.ObjectIDFromHex(todoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
			return
		}
		todos := mongoClient.GetCollection(database.CollectionNames.Todos)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = todos.DeleteOne(ctx, bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
	}
}

// CompleteTodoHandler marks a todo as complete
func CompleteTodoHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		todoID := c.Param("id")
		if todoID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Todo ID required"})
			return
		}
		objID, err := primitive.ObjectIDFromHex(todoID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
			return
		}
		todos := mongoClient.GetCollection(database.CollectionNames.Todos)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Find the todo first to get current completion status
		var todo models.Todo
		err = todos.FindOne(ctx, bson.M{"_id": objID}).Decode(&todo)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Toggle completion status
		newCompleted := !todo.Completed
		var completedAt *time.Time
		if newCompleted {
			now := time.Now()
			completedAt = &now
		}

		// Update the todo
		update := bson.M{
			"$set": bson.M{
				"completed":    newCompleted,
				"completed_at": completedAt,
				"updated_at":   time.Now(),
			},
		}

		_, err = todos.UpdateOne(ctx, bson.M{"_id": objID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Todo completion status updated"})
	}
}
