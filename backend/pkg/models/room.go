package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Room represents a study room in the system
type Room struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name            string             `bson:"name" json:"name"`
	Description     string             `bson:"description" json:"description"`
	Type            string             `bson:"type" json:"type"` // 'shared', 'private', etc.
	CreatorID       string             `bson:"creator_id" json:"creatorId"`
	Participants    []string           `bson:"participants" json:"participants"` // User IDs
	MaxParticipants int                `bson:"max_participants" json:"maxParticipants"`
	Materials       []string           `bson:"materials" json:"materials"` // Material IDs
	Todos           []string           `bson:"todos" json:"todos"`         // Todo IDs
	Notes           []string           `bson:"notes" json:"notes"`         // Note IDs
	InvitationCode  string             `bson:"invitation_code" json:"invitationCode"`
	CreatedAt       time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updatedAt"`
	LastActivityAt  time.Time          `bson:"last_activity_at" json:"lastActivityAt"`
	IsActive        bool               `bson:"is_active" json:"isActive"`
}

// RoomForResponse represents a room object for API responses
type RoomForResponse struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Type             string            `json:"type"`
	CreatorID        string            `json:"creatorId"`
	CreatorUsername  string            `json:"creatorUsername"`
	ParticipantCount int               `json:"participantCount"`
	MaxParticipants  int               `json:"maxParticipants"`
	MaterialsCount   int               `json:"materialsCount"`
	TodosCount       int               `json:"todosCount"`
	NotesCount       int               `json:"notesCount"`
	InvitationCode   string            `json:"invitationCode,omitempty"`
	Participants     []ParticipantInfo `json:"participants"`
	CreatedAt        time.Time         `json:"createdAt"`
	LastActivityAt   time.Time         `json:"lastActivityAt"`
	IsActive         bool              `json:"isActive"`
}

// ParticipantInfo represents participant information for room responses
type ParticipantInfo struct {
	UserID    string `json:"userId"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatarUrl,omitempty"`
	IsOnline  bool   `json:"isOnline"`
}

// CreateRoomRequest represents the create room request body
type CreateRoomRequest struct {
	Name            string `json:"name" binding:"required,min=3,max=100"`
	Description     string `json:"description" binding:"max=500"`
	MaxParticipants int    `json:"maxParticipants" binding:"required,min=1,max=50"`
}

// UpdateRoomRequest represents the update room request body
type UpdateRoomRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	MaxParticipants int    `json:"maxParticipants"`
}

// AddParticipantRequest represents the add participant request body
type AddParticipantRequest struct {
	UserID string `json:"userId" binding:"required"`
}

// JoinByCodeRequest represents the join by code request body
type JoinByCodeRequest struct {
	InvitationCode string `json:"invitationCode" binding:"required"`
}

// ToResponse converts a Room to a RoomForResponse
func (r *Room) ToResponse() RoomForResponse {
	return RoomForResponse{
		ID:               r.ID.Hex(),
		Name:             r.Name,
		Description:      r.Description,
		Type:             r.Type,
		CreatorID:        r.CreatorID,
		CreatorUsername:  "", // Will be populated by handler
		ParticipantCount: len(r.Participants),
		MaxParticipants:  r.MaxParticipants,
		MaterialsCount:   len(r.Materials),
		TodosCount:       len(r.Todos),
		NotesCount:       len(r.Notes),
		Participants:     []ParticipantInfo{}, // Will be populated by handler
		CreatedAt:        r.CreatedAt,
		LastActivityAt:   r.LastActivityAt,
		IsActive:         r.IsActive,
	}
}
