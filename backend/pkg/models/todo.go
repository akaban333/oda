package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Todo represents a todo item
type Todo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	Completed   bool               `bson:"completed" json:"completed"`
	DueDate     *time.Time         `bson:"due_date,omitempty" json:"dueDate,omitempty"`
	Priority    int                `bson:"priority" json:"priority"` // 1 (low) to 3 (high)
	RoomID      string             `bson:"room_id" json:"roomId"`
	CreatorID   string             `bson:"creator_id" json:"creatorId"`
	AssigneeIDs []string           `bson:"assignee_ids" json:"assigneeIds"` // User IDs
	Tags        []string           `bson:"tags" json:"tags"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
	CompletedAt *time.Time         `bson:"completed_at,omitempty" json:"completedAt,omitempty"`
}

// TodoForResponse represents a todo object for API responses
type TodoForResponse struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Completed     bool       `json:"completed"`
	DueDate       *time.Time `json:"dueDate,omitempty"`
	Priority      int        `json:"priority"`
	RoomID        string     `json:"roomId"`
	CreatorID     string     `json:"creatorId"`
	AssigneeIDs   []string   `json:"assigneeIds"`
	AssigneeCount int        `json:"assigneeCount"`
	Tags          []string   `json:"tags"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
}

// CreateTodoRequest represents the create todo request body
type CreateTodoRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"dueDate"`
	Priority    int        `json:"priority" binding:"min=1,max=3"`
	RoomID      string     `json:"roomId"`
	AssigneeIDs []string   `json:"assigneeIds"`
	Tags        []string   `json:"tags"`
}

// UpdateTodoRequest represents the update todo request body
type UpdateTodoRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	DueDate     *time.Time `json:"dueDate"`
	Priority    int        `json:"priority" binding:"min=1,max=3"`
	AssigneeIDs []string   `json:"assigneeIds"`
	Tags        []string   `json:"tags"`
}

// AssignTodoRequest represents the assign todo request body
type AssignTodoRequest struct {
	UserID string `json:"userId" binding:"required"`
}

// ToResponse converts a Todo to a TodoForResponse
func (t *Todo) ToResponse() TodoForResponse {
	return TodoForResponse{
		ID:            t.ID.Hex(),
		Title:         t.Title,
		Description:   t.Description,
		Completed:     t.Completed,
		DueDate:       t.DueDate,
		Priority:      t.Priority,
		RoomID:        t.RoomID,
		CreatorID:     t.CreatorID,
		AssigneeIDs:   t.AssigneeIDs,
		AssigneeCount: len(t.AssigneeIDs),
		Tags:          t.Tags,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
		CompletedAt:   t.CompletedAt,
	}
}
