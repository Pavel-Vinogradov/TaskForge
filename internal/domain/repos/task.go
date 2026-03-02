package repos

import (
	"TaskForge/internal/domain/entity"
	"context"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task entity.Task) (entity.Task, error)
	ListTasks(ctx context.Context, filters TaskFilters) ([]entity.Task, int64, error)
}

type TaskFilters struct {
	TeamID     *int
	Status     *string
	AssigneeID *int
	Page       int
	Limit      int
}
