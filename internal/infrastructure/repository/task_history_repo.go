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

func (r *taskHistoryRepository) GetTaskHistory(ctx context.Context, taskID int) ([]entity.TaskHistory, error) {
	query := sq.Select("*").
		From("task_history").
		Where(sq.Eq{"task_id": taskID}).
		OrderBy("changed_at DESC").
		RunWith(r.db).
		PlaceholderFormat(sq.Question)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query task history: %w", err)
	}
	defer rows.Close()

	var history []entity.TaskHistory
	for rows.Next() {
		var h entity.TaskHistory
		err := rows.Scan(
			&h.Id, &h.TaskID, &h.ChangedBy, &h.Field,
			&h.OldValue, &h.NewValue, &h.ChangedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task history: %w", err)
		}
		history = append(history, h)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating task history: %w", err)
	}

	return history, nil
}
