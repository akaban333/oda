package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MaterialType defines the type of material
type MaterialType string

const (
	// MaterialFile is a file uploaded to the system
	MaterialFile MaterialType = "file"
	// MaterialLink is a link to an external resource
	MaterialLink MaterialType = "link"
)

// Material represents a study material
type Material struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Description  string             `bson:"description" json:"description"`
	Type         MaterialType       `bson:"type" json:"type"`
	URL          string             `bson:"url" json:"url"`
	FileMetadata *FileMetadata      `bson:"file_metadata,omitempty" json:"fileMetadata,omitempty"`
	RoomID       string             `bson:"room_id" json:"roomId"`
	OwnerID      string             `bson:"owner_id" json:"ownerId"`
	Tags         []string           `bson:"tags" json:"tags"`
	ViewCount    int                `bson:"view_count" json:"viewCount"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

// FileMetadata contains metadata for uploaded files
type FileMetadata struct {
	OriginalName string `bson:"original_name" json:"originalName"`
	Size         int64  `bson:"size" json:"size"` // In bytes
	ContentType  string `bson:"content_type" json:"contentType"`
	Extension    string `bson:"extension" json:"extension"`
}

// MaterialForResponse represents a material object for API responses
type MaterialForResponse struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Type         MaterialType  `json:"type"`
	URL          string        `json:"url,omitempty"`
	FileMetadata *FileMetadata `json:"fileMetadata,omitempty"`
	RoomID       string        `json:"roomId"`
	OwnerID      string        `json:"ownerId"`
	Tags         []string      `json:"tags"`
	ViewCount    int           `json:"viewCount"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

// CreateMaterialRequest represents the create material request body
type CreateMaterialRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"max=500"`
	RoomID      string `json:"roomId"`
	FileType    string `json:"fileType" binding:"required"`
	FileURL     string `json:"fileUrl" binding:"required"`
	FileSize    int64  `json:"fileSize"`
}

// UpdateMaterialRequest represents the update material request body
type UpdateMaterialRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	FileURL     string `json:"fileUrl"`
	FileSize    int64  `json:"fileSize"`
}

// ToResponse converts a Material to a MaterialForResponse
func (m *Material) ToResponse() MaterialForResponse {
	return MaterialForResponse{
		ID:           m.ID.Hex(),
		Name:         m.Name,
		Description:  m.Description,
		Type:         m.Type,
		URL:          m.URL,
		FileMetadata: m.FileMetadata,
		RoomID:       m.RoomID,
		OwnerID:      m.OwnerID,
		Tags:         m.Tags,
		ViewCount:    m.ViewCount,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
