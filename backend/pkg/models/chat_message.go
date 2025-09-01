package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ChatMessage represents a chat message in a room
type ChatMessage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID    string             `bson:"room_id" json:"roomId"`
	UserID    string             `bson:"user_id" json:"userId"`
	Username  string             `bson:"username" json:"username"`
	Content   string             `bson:"content" json:"content"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}

// ChatMessageForResponse represents a chat message for API responses
type ChatMessageForResponse struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"roomId"`
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"createdAt"`
}

// ToResponse converts a ChatMessage to a ChatMessageForResponse
func (cm *ChatMessage) ToResponse() ChatMessageForResponse {
	return ChatMessageForResponse{
		ID:        cm.ID.Hex(),
		RoomID:    cm.RoomID,
		UserID:    cm.UserID,
		Username:  cm.Username,
		Content:   cm.Content,
		Timestamp: cm.Timestamp,
		CreatedAt: cm.CreatedAt,
	}
}
