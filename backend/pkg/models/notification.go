package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Notification represents a notification in the system
type Notification struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserID    string                 `bson:"user_id" json:"userId"`
	Type      string                 `bson:"type" json:"type"`
	Title     string                 `bson:"title" json:"title"`
	Message   string                 `bson:"message" json:"message"`
	TargetID  string                 `bson:"target_id,omitempty" json:"targetId,omitempty"`
	IsRead    bool                   `bson:"is_read" json:"isRead"`
	CreatedAt time.Time              `bson:"created_at" json:"createdAt"`
	Data      map[string]interface{} `bson:"data,omitempty" json:"data,omitempty"`
}

// NotificationForResponse represents a notification object for API responses
type NotificationForResponse struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Message   string                 `json:"message"`
	TargetID  string                 `json:"targetId,omitempty"`
	IsRead    bool                   `json:"isRead"`
	CreatedAt time.Time              `json:"createdAt"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// Notification types
const (
	NotificationTypeFriendRequest  = "friend_request"
	NotificationTypeFriendAccepted = "friend_accepted"
	NotificationTypePostLike       = "post_like"
	NotificationTypeCommentLike    = "comment_like"
	NotificationTypePostComment    = "post_comment"
	NotificationTypeCommentComment = "comment_comment"
	NotificationTypeRoomInvitation = "room_invitation"
	NotificationTypeRoomJoined     = "room_joined"
	NotificationTypeXPLevelUp      = "xp_level_up"
	NotificationTypeSystem         = "system"
)

// ToResponse converts a Notification to NotificationForResponse
func (n *Notification) ToResponse() NotificationForResponse {
	return NotificationForResponse{
		ID:        n.ID.Hex(),
		Type:      n.Type,
		Title:     n.Title,
		Message:   n.Message,
		TargetID:  n.TargetID,
		IsRead:    n.IsRead,
		CreatedAt: n.CreatedAt,
		Data:      n.Data,
	}
}

// CreateFriendRequestNotification creates a friend request notification
func CreateFriendRequestNotification(userID, requesterID, requesterUsername string) Notification {
	return Notification{
		UserID:    userID,
		Type:      NotificationTypeFriendRequest,
		Title:     "New Friend Request",
		Message:   requesterUsername + " sent you a friend request",
		TargetID:  requesterID,
		IsRead:    false,
		CreatedAt: time.Now(),
		Data: map[string]interface{}{
			"requesterUsername": requesterUsername,
			"requesterID":       requesterID,
		},
	}
}

// CreateFriendAcceptedNotification creates a friend accepted notification
func CreateFriendAcceptedNotification(userID, accepterID, accepterUsername string) Notification {
	return Notification{
		UserID:    userID,
		Type:      NotificationTypeFriendAccepted,
		Title:     "Friend Request Accepted",
		Message:   accepterUsername + " accepted your friend request",
		TargetID:  accepterID,
		IsRead:    false,
		CreatedAt: time.Now(),
		Data: map[string]interface{}{
			"accepterUsername": accepterUsername,
			"accepterID":       accepterID,
		},
	}
}

// CreatePostLikeNotification creates a post like notification
func CreatePostLikeNotification(userID, likerID, likerUsername, postID, postContent string) Notification {
	return Notification{
		UserID:    userID,
		Type:      NotificationTypePostLike,
		Title:     "New Like",
		Message:   likerUsername + " liked your post",
		TargetID:  postID,
		IsRead:    false,
		CreatedAt: time.Now(),
		Data: map[string]interface{}{
			"likerUsername": likerUsername,
			"likerID":       likerID,
			"postID":        postID,
			"postContent":   postContent,
		},
	}
}

// CreateCommentLikeNotification creates a comment like notification
func CreateCommentLikeNotification(userID, likerID, likerUsername, commentID, commentContent string) Notification {
	return Notification{
		UserID:    userID,
		Type:      NotificationTypeCommentLike,
		Title:     "New Like",
		Message:   likerUsername + " liked your comment",
		TargetID:  commentID,
		IsRead:    false,
		CreatedAt: time.Now(),
		Data: map[string]interface{}{
			"likerUsername":  likerUsername,
			"likerID":        likerID,
			"commentID":      commentID,
			"commentContent": commentContent,
		},
	}
}

// CreatePostCommentNotification creates a post comment notification
func CreatePostCommentNotification(userID, commenterID, commenterUsername, postID string) Notification {
	return Notification{
		UserID:    userID,
		Type:      NotificationTypePostComment,
		Title:     "New Comment",
		Message:   commenterUsername + " commented on your post",
		TargetID:  postID,
		IsRead:    false,
		CreatedAt: time.Now(),
		Data: map[string]interface{}{
			"commenterUsername": commenterUsername,
			"commenterID":       commenterID,
			"postID":            postID,
		},
	}
}

// CreateRoomInvitationNotification creates a room invitation notification
func CreateRoomInvitationNotification(userID, inviterID, inviterUsername, roomID, roomName string) Notification {
	return Notification{
		UserID:    userID,
		Type:      NotificationTypeRoomInvitation,
		Title:     "Room Invitation",
		Message:   inviterUsername + " invited you to join " + roomName,
		TargetID:  roomID,
		IsRead:    false,
		CreatedAt: time.Now(),
		Data: map[string]interface{}{
			"inviterUsername": inviterUsername,
			"inviterID":       inviterID,
			"roomID":          roomID,
			"roomName":        roomName,
		},
	}
}

// CreateXPLevelUpNotification creates an XP level up notification
func CreateXPLevelUpNotification(userID string, newLevel int) Notification {
	return Notification{
		UserID:    userID,
		Type:      NotificationTypeXPLevelUp,
		Title:     "Level Up!",
		Message:   "Congratulations! You've reached level " + string(newLevel),
		IsRead:    false,
		CreatedAt: time.Now(),
		Data: map[string]interface{}{
			"newLevel": newLevel,
		},
	}
}
