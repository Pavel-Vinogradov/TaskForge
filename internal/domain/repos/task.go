package repos

import (
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/interfaces/task"
	"context"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task entity.Task) (entity.Task, error)
	GetTaskByID(ctx context.Context, taskID int) (entity.Task, error)
	UpdateTask(ctx context.Context, task entity.Task) (entity.Task, error)
	ListTasks(ctx context.Context, filters task.TaskFilters) ([]entity.Task, int64, error)
}
