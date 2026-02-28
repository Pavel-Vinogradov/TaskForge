package task

import (
	"TaskForge/internal/domain/entity"
	"context"
)

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=100" example:"Task"`
	Description string `json:"description" binding:"min=3" example:"Description"`
	TeamID      int    `json:"team_id" binding:"required" example:"1"`
	AssigneeID  int    `json:"assignee_id" binding:"required" example:"1"`
}

type UseCaseTask interface {
	CreateTask(ctx context.Context, req CreateTaskRequest) (*entity.Task, error)
	ListTask(ctx context.Context) ([]entity.Task, error)
	UpdateTask(ctx context.Context) (*entity.Task, error)
	HistoryTask(ctx context.Context) (interface{}, error)
}
