package material

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
	"github.com/studyplatform/backend/pkg/storage"
)

// ListMaterialsHandler returns all materials for the authenticated user
func ListMaterialsHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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
		materials := mongoClient.GetCollection(database.CollectionNames.Materials)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		filter := bson.M{"owner_id": userIDStr}
		cursor, err := materials.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer cursor.Close(ctx)
		var materialList []models.MaterialForResponse
		for cursor.Next(ctx) {
			var material models.Material
			if err := cursor.Decode(&material); err == nil {
				materialList = append(materialList, material.ToResponse())
			}
		}
		c.JSON(http.StatusOK, gin.H{"materials": materialList})
	}
}

// CreateMaterialHandler creates a new study material
func CreateMaterialHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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
		var req models.CreateMaterialRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		materials := mongoClient.GetCollection(database.CollectionNames.Materials)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		materialMap := bson.M{
			"name":        req.Name,
			"description": req.Description,
			"owner_id":    userIDStr,
			"room_id":     req.RoomID,
			"file_type":   req.FileType,
			"file_url":    req.FileURL,
			"file_size":   req.FileSize,
			"created_at":  time.Now(),
			"shared_with": []bson.M{},
		}
		res, err := materials.InsertOne(ctx, materialMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create material"})
			return
		}

		// Get the inserted ID
		insertedID, ok := res.InsertedID.(primitive.ObjectID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inserted ID"})
			return
		}

		// Fetch the created material to return it
		var material models.Material
		err = materials.FindOne(ctx, bson.M{"_id": insertedID}).Decode(&material)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created material"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"material": material.ToResponse()})
	}
}

// GetMaterialHandler returns a specific material by ID
func GetMaterialHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		materialID := c.Param("id")
		if materialID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Material ID required"})
			return
		}
		objID, err := primitive.ObjectIDFromHex(materialID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
			return
		}
		materials := mongoClient.GetCollection(database.CollectionNames.Materials)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var material models.Material
		err = materials.FindOne(ctx, bson.M{"_id": objID}).Decode(&material)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Material not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"material": material.ToResponse()})
	}
}

// UpdateMaterialHandler updates a material
func UpdateMaterialHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		materialID := c.Param("id")
		if materialID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Material ID required"})
			return
		}
		materialObjID, err := primitive.ObjectIDFromHex(materialID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
			return
		}
		var req models.UpdateMaterialRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}
		materials := mongoClient.GetCollection(database.CollectionNames.Materials)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		updateFields := bson.M{"updated_at": time.Now()}
		if req.Name != "" {
			updateFields["name"] = req.Name
		}
		if req.Description != "" {
			updateFields["description"] = req.Description
		}
		if req.FileURL != "" {
			updateFields["file_url"] = req.FileURL
		}
		if req.FileSize > 0 {
			updateFields["file_size"] = req.FileSize
		}
		update := bson.M{"$set": updateFields}
		res, err := materials.UpdateOne(ctx, bson.M{"_id": materialObjID, "owner_id": userIDStr}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update material"})
			return
		}
		if res.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Material not found or you don't have permission"})
			return
		}
		var material models.Material
		err = materials.FindOne(ctx, bson.M{"_id": materialObjID}).Decode(&material)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated material"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"material": material.ToResponse()})
	}
}

// DeleteMaterialHandler deletes a material
func DeleteMaterialHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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
		materialID := c.Param("id")
		if materialID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Material ID required"})
			return
		}
		materialObjID, err := primitive.ObjectIDFromHex(materialID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
			return
		}
		materials := mongoClient.GetCollection(database.CollectionNames.Materials)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		res, err := materials.DeleteOne(ctx, bson.M{"_id": materialObjID, "owner_id": userIDStr})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete material"})
			return
		}
		if res.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Material not found or you don't have permission"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Material deleted successfully"})
	}
}

// GenerateUploadURLHandler generates a presigned URL for file upload
func GenerateUploadURLHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		var req struct {
			FileName    string `json:"fileName" binding:"required"`
			FileType    string `json:"fileType" binding:"required"`
			ContentType string `json:"contentType" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		// Initialize MinIO client
		minioClient, err := storage.NewMinioClient()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize storage client"})
			return
		}

		// Generate unique object name
		objectName := generateUniqueFileName(req.FileName)

		// Generate presigned URL for upload
		uploadURL, err := minioClient.GetPresignedUploadURL("materials", objectName, req.ContentType, 15*time.Minute)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate upload URL"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"uploadURL":  uploadURL,
			"objectName": objectName,
			"expiresIn":  "15 minutes",
		})
	}
}

// ConfirmUploadHandler confirms the file upload and creates a material record
func ConfirmUploadHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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
		objID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var req struct {
			ObjectName  string `json:"objectName" binding:"required"`
			Name        string `json:"name" binding:"required"`
			Description string `json:"description"`
			FileType    string `json:"fileType" binding:"required"`
			FileSize    int64  `json:"fileSize" binding:"required"`
			RoomID      string `json:"roomId"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		// Initialize MinIO client to verify file exists
		minioClient, err := storage.NewMinioClient()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize storage client"})
			return
		}

		// Verify file exists in storage
		fileExists, err := minioClient.FileExists("materials", req.ObjectName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify file upload"})
			return
		}
		if !fileExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File not found in storage"})
			return
		}

		// Generate file URL
		fileURL, err := minioClient.GetFileURL("materials", req.ObjectName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate file URL"})
			return
		}

		materials := mongoClient.GetCollection(database.CollectionNames.Materials)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		materialMap := bson.M{
			"name":        req.Name,
			"description": req.Description,
			"owner_id":    objID,
			"file_type":   req.FileType,
			"file_url":    fileURL,
			"file_size":   req.FileSize,
			"object_name": req.ObjectName,
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		}

		if req.RoomID != "" {
			roomObjID, err := primitive.ObjectIDFromHex(req.RoomID)
			if err == nil {
				materialMap["room_id"] = roomObjID
			}
		}

		res, err := materials.InsertOne(ctx, materialMap)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create material record"})
			return
		}

		var material models.Material
		if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
			material.ID = oid
		}

		err = materials.FindOne(ctx, bson.M{"_id": material.ID}).Decode(&material)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created material"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"material": material.ToResponse()})
	}
}

// ShareMaterialHandler shares a material with another user
func ShareMaterialHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
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

		materialID := c.Param("id")
		if materialID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Material ID required"})
			return
		}
		materialObjID, err := primitive.ObjectIDFromHex(materialID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid material ID"})
			return
		}

		var req struct {
			UserID     string `json:"userId" binding:"required"`
			Permission string `json:"permission" binding:"required,oneof=view edit"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
			return
		}

		materials := mongoClient.GetCollection(database.CollectionNames.Materials)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Verify user owns the material
		var material models.Material
		err = materials.FindOne(ctx, bson.M{"_id": materialObjID, "owner_id": userIDStr}).Decode(&material)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Material not found or you don't have permission"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Add or update share permission
		shareEntry := bson.M{
			"user_id":    req.UserID,
			"permission": req.Permission,
			"shared_at":  time.Now(),
		}

		update := bson.M{
			"$addToSet": bson.M{"shared_with": shareEntry},
			"$set":      bson.M{"updated_at": time.Now()},
		}

		_, err = materials.UpdateOne(ctx, bson.M{"_id": materialObjID}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to share material"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Material shared successfully"})
	}
}

// GetRoomMaterialsHandler retrieves materials for a specific room
func GetRoomMaterialsHandler(mongoClient *database.MongoClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("id")
		if roomID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID required"})
			return
		}

		materials := mongoClient.GetCollection(database.CollectionNames.Materials)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Get materials for the specific room
		filter := bson.M{"room_id": roomID}
		cursor, err := materials.Find(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch room materials"})
			return
		}
		defer cursor.Close(ctx)

		var materialsList []models.Material
		if err = cursor.All(ctx, &materialsList); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode room materials"})
			return
		}

		// Convert to response format
		var responseMaterials []models.MaterialForResponse
		for _, material := range materialsList {
			responseMaterials = append(responseMaterials, material.ToResponse())
		}

		c.JSON(http.StatusOK, gin.H{"materials": responseMaterials})
	}
}

// Helper function to generate unique file names
func generateUniqueFileName(originalName string) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d_%s", timestamp, originalName)
}
