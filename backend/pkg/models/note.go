package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Note represents a study note
type Note struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title        string             `bson:"title" json:"title"`
	Content      string             `bson:"content" json:"content"`
	RoomID       string             `bson:"room_id" json:"roomId"`
	CreatorID    string             `bson:"creator_id" json:"creatorId"`
	Tags         []string           `bson:"tags" json:"tags"`
	SharedWith   []string           `bson:"shared_with" json:"sharedWith"` // User IDs
	IsPublic     bool               `bson:"is_public" json:"isPublic"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
	LastEditedBy string             `bson:"last_edited_by" json:"lastEditedBy"`
}

// NoteRevision represents a specific revision of a note
type NoteRevision struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	NoteID    string             `bson:"note_id" json:"noteId"`
	Content   string             `bson:"content" json:"content"`
	EditorID  string             `bson:"editor_id" json:"editorId"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}

// NoteForResponse represents a note object for API responses
type NoteForResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	RoomID       string    `json:"roomId"`
	CreatorID    string    `json:"creatorId"`
	Tags         []string  `json:"tags"`
	SharedWith   []string  `json:"sharedWith,omitempty"`
	SharedCount  int       `json:"sharedCount"`
	IsPublic     bool      `json:"isPublic"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	LastEditedBy string    `json:"lastEditedBy"`
}

// CreateNoteRequest represents the create note request body
type CreateNoteRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=1000"`
	RoomID   string `json:"roomId"`
	IsShared bool   `json:"isShared"`
}

// UpdateNoteRequest represents the update note request body
type UpdateNoteRequest struct {
	Content  string `json:"content"`
	IsShared *bool  `json:"isShared"`
}

// ShareNoteRequest represents the share note request body
type ShareNoteRequest struct {
	UserIDs []string `json:"userIds" binding:"required"`
}

// ToResponse converts a Note to a NoteForResponse
func (n *Note) ToResponse() NoteForResponse {
	return NoteForResponse{
		ID:           n.ID.Hex(),
		Title:        n.Title,
		Content:      n.Content,
		RoomID:       n.RoomID,
		CreatorID:    n.CreatorID,
		Tags:         n.Tags,
		SharedWith:   n.SharedWith,
		SharedCount:  len(n.SharedWith),
		IsPublic:     n.IsPublic,
		CreatedAt:    n.CreatedAt,
		UpdatedAt:    n.UpdatedAt,
		LastEditedBy: n.LastEditedBy,
	}
}
