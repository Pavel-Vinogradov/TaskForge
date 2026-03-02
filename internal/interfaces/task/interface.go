package task

import (
	"TaskForge/internal/domain/entity"
	"context"
	"time"
)

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=100" example:"Task"`
	Description string `json:"description" binding:"min=3" example:"Description"`
	TeamID      int    `json:"team_id" binding:"required" example:"1"`
	AssigneeID  int    `json:"assignee_id" binding:"required" example:"1"`
}

type TaskListRequest struct {
	TeamID     *int    `form:"team_id" example:"1"`
	Status     *string `form:"status" example:"todo"`
	AssigneeID *int    `form:"assignee_id" example:"1"`
	Page       int     `form:"page,default=1" example:"1"`
	Limit      int     `form:"limit,default=10" example:"10"`
}

type TaskListResult struct {
	Tasks []entity.Task `json:"tasks"`
	Total int64         `json:"total"`
}

type ResponseTask struct {
	ID          int       `json:"id" example:"1"`
	Title       string    `json:"title" example:"Task"`
	Description string    `json:"description" example:"Description"`
	Status      string    `json:"status" example:"todo"`
	TeamID      int       `json:"team_id" example:"1"`
	CreatedBy   int       `json:"created_by" example:"1"`
	AssigneeID  int       `json:"assignee_id" example:"1"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

type TaskFilters struct {
	TeamID     *int
	Status     *string
	AssigneeID *int
	Page       int
	Limit      int
}
type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
	AssigneeID  *int    `json:"assignee_id,omitempty"`
}

type UseCaseTask interface {
	CreateTask(ctx context.Context, req CreateTaskRequest) (ResponseTask, error)
	ListTask(ctx context.Context, req TaskListRequest) (TaskListResult, error)
	UpdateTask(ctx context.Context, taskID int, req UpdateTaskRequest) (ResponseTask, error)
	HistoryTask(ctx context.Context, taskID int) ([]entity.TaskHistory, error)
}
