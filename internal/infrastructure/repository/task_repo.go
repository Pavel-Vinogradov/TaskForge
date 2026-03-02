package repository

import (
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/domain/repos"
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) repos.TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) CreateTask(ctx context.Context, task entity.Task) (entity.Task, error) {
	query := sq.Insert("tasks").
		Columns("title", "description", "status", "assignee_id", "team_id", "created_by", "created_at").
		Values(task.Title, task.Description, task.Status, task.AssigneeID, task.TeamID, task.CreatedBy, task.CreatedAt).
		RunWith(r.db).
		PlaceholderFormat(sq.Question)

	res, err := query.ExecContext(ctx)
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to insert task: %w", err)

	}
	id, err := res.LastInsertId()
	if err != nil {
		return entity.Task{}, fmt.Errorf("failed to get last insert id: %w", err)
	}

	task.Id = int(id)
	return task, nil

}

func (r *taskRepository) ListTasks(ctx context.Context, filters repos.TaskFilters) ([]entity.Task, int64, error) {
	countQuery := sq.Select("COUNT(*)").From("tasks")

	if filters.TeamID != nil {
		countQuery = countQuery.Where("team_id = ?", *filters.TeamID)
	}
	if filters.Status != nil {
		countQuery = countQuery.Where("status = ?", *filters.Status)
	}
	if filters.AssigneeID != nil {
		countQuery = countQuery.Where("assignee_id = ?", *filters.AssigneeID)
	}

	var total int64
	err := countQuery.RunWith(r.db).PlaceholderFormat(sq.Question).QueryRowContext(ctx).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	query := sq.Select("*").
		From("tasks").
		RunWith(r.db).
		PlaceholderFormat(sq.Question)

	if filters.TeamID != nil {
		query = query.Where("team_id = ?", *filters.TeamID)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.AssigneeID != nil {
		query = query.Where("assignee_id = ?", *filters.AssigneeID)
	}

	offset := (filters.Page - 1) * filters.Limit
	query = query.Limit(uint64(uint(filters.Limit))).Offset(uint64(offset))

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var tasks []entity.Task
	for rows.Next() {
		var task entity.Task
		err := rows.Scan(
			&task.Id, &task.Title, &task.Description, &task.Status,
			&task.TeamID, &task.CreatedBy, &task.AssigneeID,
			&task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, total, nil
}
