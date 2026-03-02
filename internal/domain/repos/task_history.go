package repos

import (
	"TaskForge/internal/domain/entity"
	"context"
)

type TaskHistoryRepository interface {
	CreateHistory(ctx context.Context, history entity.TaskHistory) error
}
