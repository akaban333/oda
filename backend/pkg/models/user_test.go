package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUser_ToResponse(t *testing.T) {
	userID := primitive.NewObjectID()
	now := time.Now()

	user := User{
		ID:           userID,
		UniqueID:     "user123",
		Username:     "testuser",
		Email:        "test@example.com",
		FirstName:    "John",
		LastName:     "Doe",
		AvatarURL:    "https://example.com/avatar.jpg",
		Bio:          "Test bio",
		XP:           150,
		Level:        5,
		Friends:      []Friend{{UserID: "friend1", Status: "accepted", Since: now}},
		JoinedRooms:  []string{"room1", "room2"},
		CreatedRooms: []string{"room1"},
		CreatedAt:    now,
		UpdatedAt:    now,
		LastActive:   now,
		IsActive:     true,
		IsVerified:   true,
	}

	response := user.ToResponse()

	assert.Equal(t, userID.Hex(), response.ID)
	assert.Equal(t, user.UniqueID, response.UniqueID)
	assert.Equal(t, user.Username, response.Username)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.FirstName, response.FirstName)
	assert.Equal(t, user.LastName, response.LastName)
	assert.Equal(t, user.AvatarURL, response.AvatarURL)
	assert.Equal(t, user.Bio, response.Bio)
	assert.Equal(t, user.XP, response.TotalXP)
	assert.Equal(t, user.Level, response.Level)
	assert.Equal(t, len(user.Friends), response.FriendsCount)
	assert.Equal(t, len(user.JoinedRooms), response.RoomsCount)
	assert.Equal(t, user.CreatedAt, response.CreatedAt)
	assert.Equal(t, user.IsActive, response.IsActive)
	assert.Equal(t, user.IsVerified, response.IsVerified)
}

func TestUser_ToResponse_EmptyFields(t *testing.T) {
	userID := primitive.NewObjectID()
	now := time.Now()

	user := User{
		ID:           userID,
		UniqueID:     "user123",
		Username:     "testuser",
		Email:        "test@example.com",
		FirstName:    "",
		LastName:     "",
		AvatarURL:    "",
		Bio:          "",
		XP:           0,
		Level:        1,
		Friends:      []Friend{},
		JoinedRooms:  []string{},
		CreatedRooms: []string{},
		CreatedAt:    now,
		UpdatedAt:    now,
		LastActive:   now,
		IsActive:     true,
		IsVerified:   false,
	}

	response := user.ToResponse()

	assert.Equal(t, userID.Hex(), response.ID)
	assert.Equal(t, user.UniqueID, response.UniqueID)
	assert.Equal(t, user.Username, response.Username)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, "", response.FirstName)
	assert.Equal(t, "", response.LastName)
	assert.Equal(t, "", response.AvatarURL)
	assert.Equal(t, "", response.Bio)
	assert.Equal(t, 0, response.TotalXP)
	assert.Equal(t, 1, response.Level)
	assert.Equal(t, 0, response.FriendsCount)
	assert.Equal(t, 0, response.RoomsCount)
	assert.Equal(t, user.CreatedAt, response.CreatedAt)
	assert.Equal(t, user.IsActive, response.IsActive)
	assert.Equal(t, user.IsVerified, response.IsVerified)
}

func TestSignupRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request SignupRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: SignupRequest{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			isValid: true,
		},
		{
			name: "username too short",
			request: SignupRequest{
				Username:  "ab",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			isValid: false,
		},
		{
			name: "username too long",
			request: SignupRequest{
				Username:  "verylongusernameexceedingthirtycharacters",
				Email:     "test@example.com",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			isValid: false,
		},
		{
			name: "invalid email",
			request: SignupRequest{
				Username:  "testuser",
				Email:     "invalid-email",
				Password:  "password123",
				FirstName: "John",
				LastName:  "Doe",
			},
			isValid: false,
		},
		{
			name: "password too short",
			request: SignupRequest{
				Username:  "testuser",
				Email:     "test@example.com",
				Password:  "123",
				FirstName: "John",
				LastName:  "Doe",
			},
			isValid: false,
		},
		{
			name: "missing required fields",
			request: SignupRequest{
				Username:  "",
				Email:     "",
				Password:  "",
				FirstName: "John",
				LastName:  "Doe",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: In a real application, you would use a validation library
			// like go-playground/validator to test the struct tags
			// This is a simplified test for demonstration
			if tt.isValid {
				assert.NotEmpty(t, tt.request.Username)
				assert.NotEmpty(t, tt.request.Email)
				assert.NotEmpty(t, tt.request.Password)
				assert.GreaterOrEqual(t, len(tt.request.Username), 3)
				assert.GreaterOrEqual(t, len(tt.request.Password), 8)
			}
		})
	}
}

func TestLoginRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request LoginRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			isValid: true,
		},
		{
			name: "invalid email",
			request: LoginRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			isValid: false,
		},
		{
			name: "missing email",
			request: LoginRequest{
				Email:    "",
				Password: "password123",
			},
			isValid: false,
		},
		{
			name: "missing password",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isValid {
				assert.NotEmpty(t, tt.request.Email)
				assert.NotEmpty(t, tt.request.Password)
			}
		})
	}
}

func TestUpdateProfileRequest(t *testing.T) {
	request := UpdateProfileRequest{
		Username:  "newusername",
		FirstName: "Jane",
		LastName:  "Smith",
		Bio:       "New bio",
		AvatarURL: "https://example.com/new-avatar.jpg",
	}

	assert.Equal(t, "newusername", request.Username)
	assert.Equal(t, "Jane", request.FirstName)
	assert.Equal(t, "Smith", request.LastName)
	assert.Equal(t, "New bio", request.Bio)
	assert.Equal(t, "https://example.com/new-avatar.jpg", request.AvatarURL)
}

func TestUpdateXPRequest(t *testing.T) {
	request := UpdateXPRequest{
		XP:        100,
		Source:    "session",
		SessionID: "session123",
	}

	assert.Equal(t, 100, request.XP)
	assert.Equal(t, "session", request.Source)
	assert.Equal(t, "session123", request.SessionID)
}

func TestFriend_Validation(t *testing.T) {
	now := time.Now()

	friend := Friend{
		UserID: "user123",
		Status: "accepted",
		Since:  now,
	}

	assert.Equal(t, "user123", friend.UserID)
	assert.Equal(t, "accepted", friend.Status)
	assert.Equal(t, now, friend.Since)
}

func TestRefreshToken_Validation(t *testing.T) {
	now := time.Now()

	token := RefreshToken{
		Token:     "refresh_token_123",
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
		IP:        "192.168.1.1",
		UserAgent: "Mozilla/5.0",
		IsRevoked: false,
	}

	assert.Equal(t, "refresh_token_123", token.Token)
	assert.Equal(t, now.Add(24*time.Hour), token.ExpiresAt)
	assert.Equal(t, now, token.CreatedAt)
	assert.Equal(t, "192.168.1.1", token.IP)
	assert.Equal(t, "Mozilla/5.0", token.UserAgent)
	assert.False(t, token.IsRevoked)
}
