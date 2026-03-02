package task_history

import (
	"TaskForge/internal/domain/entity"
	"context"
)

type TaskObserver interface {
	OnTaskCreated(ctx context.Context, task entity.Task, userID int)
	OnTaskUpdated(ctx context.Context, oldTask, newTask entity.Task, userID int)
}
