package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post represents a social media post
type Post struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Content   string             `bson:"content" json:"content"`
	AuthorID  string             `bson:"author_id" json:"authorId"`
	Likes     []string           `bson:"likes" json:"likes"` // User IDs who liked
	Comments  []Comment          `bson:"comments" json:"comments"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID    string             `bson:"post_id" json:"postId"`
	Content   string             `bson:"content" json:"content"`
	AuthorID  string             `bson:"author_id" json:"authorId"`
	Likes     []string           `bson:"likes" json:"likes"` // User IDs who liked
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
}

// PostForResponse represents a post object for API responses
type PostForResponse struct {
	ID             string               `json:"id"`
	Content        string               `json:"content"`
	AuthorID       string               `json:"authorId"`
	AuthorUniqueID string               `json:"authorUniqueId"`
	AuthorName     string               `json:"authorName"`
	LikesCount     int                  `json:"likesCount"`
	IsLiked        bool                 `json:"isLiked"`
	Comments       []CommentForResponse `json:"comments"`
	CommentsCount  int                  `json:"commentsCount"`
	CreatedAt      time.Time            `json:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt"`
}

// CommentForResponse represents a comment object for API responses
type CommentForResponse struct {
	ID             string    `json:"id"`
	PostID         string    `json:"postId"`
	Content        string    `json:"content"`
	AuthorID       string    `json:"authorId"`
	AuthorUniqueID string    `json:"authorUniqueId"`
	AuthorName     string    `json:"authorName"`
	LikesCount     int       `json:"likesCount"`
	IsLiked        bool      `json:"isLiked"`
	CreatedAt      time.Time `json:"createdAt"`
}

// CreatePostRequest represents the create post request body
type CreatePostRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

// CreateCommentRequest represents the create comment request body
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=500"`
}

// LikePostRequest represents the like post request body
type LikePostRequest struct {
	PostID string `json:"postId" binding:"required"`
}

// LikeCommentRequest represents the like comment request body
type LikeCommentRequest struct {
	CommentID string `json:"commentId" binding:"required"`
}

// ToResponse converts a Post to a PostForResponse
func (p *Post) ToResponse(currentUserID string) PostForResponse {
	// Check if current user liked this post
	isLiked := false
	for _, likeID := range p.Likes {
		if likeID == currentUserID {
			isLiked = true
			break
		}
	}

	// Convert comments to response format
	comments := make([]CommentForResponse, len(p.Comments))
	for i, comment := range p.Comments {
		// Check if current user liked this comment
		commentIsLiked := false
		for _, likeID := range comment.Likes {
			if likeID == currentUserID {
				commentIsLiked = true
				break
			}
		}

		comments[i] = CommentForResponse{
			ID:             comment.ID.Hex(),
			PostID:         comment.PostID,
			Content:        comment.Content,
			AuthorID:       comment.AuthorID,
			AuthorUniqueID: "", // Will be populated by handler
			AuthorName:     "", // Will be populated by handler
			LikesCount:     len(comment.Likes),
			IsLiked:        commentIsLiked,
			CreatedAt:      comment.CreatedAt,
		}
	}

	return PostForResponse{
		ID:             p.ID.Hex(),
		Content:        p.Content,
		AuthorID:       p.AuthorID,
		AuthorUniqueID: "", // Will be populated by handler
		AuthorName:     "", // Will be populated by handler
		LikesCount:     len(p.Likes),
		IsLiked:        isLiked,
		Comments:       comments,
		CommentsCount:  len(comments),
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
