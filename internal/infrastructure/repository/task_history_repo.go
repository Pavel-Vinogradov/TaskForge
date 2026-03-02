package repository

import (
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/domain/repos"
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type taskHistoryRepository struct {
	db *sql.DB
}

func NewTaskHistoryRepository(db *sql.DB) repos.TaskHistoryRepository {
	return &taskHistoryRepository{db: db}
}

func (r *taskHistoryRepository) CreateHistory(ctx context.Context, history entity.TaskHistory) error {
	query := sq.Insert("task_history").
		Columns("task_id", "field_name", "old_value", "new_value", "changed_by", "changed_at").
		Values(history.TaskID, history.Field, history.OldValue, history.NewValue, history.ChangedBy, history.ChangedAt).
		RunWith(r.db).
		PlaceholderFormat(sq.Question)

	_, err := query.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert task history: %w", err)
	}

	return nil
}
