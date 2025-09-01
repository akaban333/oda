package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user in the system
type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UniqueID      string             `bson:"unique_id" json:"uniqueId"`
	Username      string             `bson:"username" json:"username"`
	Email         string             `bson:"email" json:"email"`
	PasswordHash  string             `bson:"password_hash" json:"-"`
	FirstName     string             `bson:"first_name" json:"firstName"`
	LastName      string             `bson:"last_name" json:"lastName"`
	AvatarURL     string             `bson:"avatar_url" json:"avatarUrl"`
	Bio           string             `bson:"bio" json:"bio"`
	XP            int                `bson:"xp" json:"xp"`
	Level         int                `bson:"level" json:"level"`
	Friends       []Friend           `bson:"friends" json:"friends"`
	JoinedRooms   []string           `bson:"joined_rooms" json:"joinedRooms"`
	CreatedRooms  []string           `bson:"created_rooms" json:"createdRooms"`
	RefreshTokens []RefreshToken     `bson:"refresh_tokens" json:"-"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
	LastActive    time.Time          `bson:"last_active" json:"lastActive"`
	IsActive      bool               `bson:"is_active" json:"isActive"`
	IsVerified    bool               `bson:"is_verified" json:"isVerified"`
}

// UserForResponse represents a user object for API responses
type UserForResponse struct {
	ID           string    `json:"id"`
	UniqueID     string    `json:"uniqueId"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	AvatarURL    string    `json:"avatarUrl"`
	Bio          string    `json:"bio"`
	TotalXP      int       `json:"totalXP"`
	Level        int       `json:"level"`
	FriendsCount int       `json:"friendsCount"`
	RoomsCount   int       `json:"roomsCount"`
	CreatedAt    time.Time `json:"createdAt"`
	IsActive     bool      `json:"isActive"`
	IsVerified   bool      `json:"isVerified"`
}

// RefreshToken represents a refresh token
type RefreshToken struct {
	Token     string    `bson:"token"`
	ExpiresAt time.Time `bson:"expires_at"`
	CreatedAt time.Time `bson:"created_at"`
	IP        string    `bson:"ip"`
	UserAgent string    `bson:"user_agent"`
	IsRevoked bool      `bson:"is_revoked"`
}

// Friend represents a friend relationship or request
// Status: "pending" (you sent), "requested" (you received), "accepted"
type Friend struct {
	UserID string    `bson:"user_id" json:"userId"`
	Status string    `bson:"status" json:"status"`
	Since  time.Time `bson:"since" json:"since"`
}

// SignupRequest represents the signup request body
type SignupRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=30"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateProfileRequest represents the update profile request body
type UpdateProfileRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatarUrl"`
}

// UpdateXPRequest represents the update XP request body
type UpdateXPRequest struct {
	XP        int    `json:"xp" binding:"required"`
	Source    string `json:"source" binding:"required"` // "session", "pomodoro", "activity"
	SessionID string `json:"sessionId,omitempty"`
}

// XPHistory represents an XP earning event
type XPHistory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"userId"`
	XP        int                `bson:"xp" json:"xp"`
	Source    string             `bson:"source" json:"source"`
	SessionID string             `bson:"session_id,omitempty" json:"sessionId,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}

// ChangePasswordRequest represents the change password request body
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
}

// ToResponse converts a User to a UserForResponse
func (u *User) ToResponse() UserForResponse {
	// Count accepted friends
	acceptedCount := 0
	for _, friend := range u.Friends {
		if friend.Status == "accepted" {
			acceptedCount++
		}
	}

	return UserForResponse{
		ID:           u.ID.Hex(),
		UniqueID:     u.UniqueID,
		Username:     u.Username,
		Email:        u.Email,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		AvatarURL:    u.AvatarURL,
		Bio:          u.Bio,
		TotalXP:      u.XP,
		Level:        u.Level,
		FriendsCount: acceptedCount,
		RoomsCount:   len(u.JoinedRooms),
		CreatedAt:    u.CreatedAt,
		IsActive:     u.IsActive,
		IsVerified:   u.IsVerified,
	}
}
