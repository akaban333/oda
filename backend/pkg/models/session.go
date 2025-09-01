package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Session represents a study session
type Session struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RoomID       string             `bson:"room_id" json:"roomId"`
	UserID       string             `bson:"user_id" json:"userId"`
	StartTime    time.Time          `bson:"start_time" json:"startTime"`
	EndTime      time.Time          `bson:"end_time" json:"endTime"`
	Duration     int64              `bson:"duration" json:"duration"` // In seconds
	XPEarned     int                `bson:"xp_earned" json:"xpEarned"`
	IsActive     bool               `bson:"is_active" json:"isActive"`
	InactiveTime int64              `bson:"inactive_time" json:"inactiveTime"` // In seconds
	Activities   []Activity         `bson:"activities" json:"activities"`
}

// Activity represents a user activity during a session
type Activity struct {
	Type      string    `bson:"type" json:"type"` // "material_view", "todo_complete", "note_create", etc.
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Details   string    `bson:"details" json:"details"` // Additional contextual information
}

// SessionForResponse represents a session object for API responses
type SessionForResponse struct {
	ID           string     `json:"id"`
	RoomID       string     `json:"roomId"`
	UserID       string     `json:"userId"`
	StartTime    time.Time  `json:"startTime"`
	EndTime      time.Time  `json:"endTime"`
	Duration     int64      `json:"duration"` // In seconds
	XPEarned     int        `json:"xpEarned"`
	IsActive     bool       `json:"isActive"`
	InactiveTime int64      `json:"inactiveTime"` // In seconds
	Activities   []Activity `json:"activities"`
}

// SessionSummary represents a summary of a user's sessions
type SessionSummary struct {
	TotalSessions   int   `json:"totalSessions"`
	TotalTimeSpent  int64 `json:"totalTimeSpent"` // In seconds
	TotalXPEarned   int   `json:"totalXPEarned"`
	AvgSessionTime  int64 `json:"avgSessionTime"`  // In seconds
	LongestSession  int64 `json:"longestSession"`  // In seconds
	ShortestSession int64 `json:"shortestSession"` // In seconds
}

// StartSessionRequest represents the start session request body
type StartSessionRequest struct {
	RoomID string `json:"roomId" binding:"required"`
}

// EndSessionRequest represents the end session request body
type EndSessionRequest struct {
	SessionID string `json:"sessionId" binding:"required"`
}

// ActivityRequest represents the activity request body
type ActivityRequest struct {
	Type    string `json:"type" binding:"required"`
	Details string `json:"details"`
}

// ToResponse converts a Session to a SessionForResponse
func (s *Session) ToResponse() SessionForResponse {
	return SessionForResponse{
		ID:           s.ID.Hex(),
		RoomID:       s.RoomID,
		UserID:       s.UserID,
		StartTime:    s.StartTime,
		EndTime:      s.EndTime,
		Duration:     s.Duration,
		XPEarned:     s.XPEarned,
		IsActive:     s.IsActive,
		InactiveTime: s.InactiveTime,
		Activities:   s.Activities,
	}
}
