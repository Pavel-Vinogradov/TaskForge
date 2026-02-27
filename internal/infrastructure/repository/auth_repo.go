package repository

import (
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/domain/repos"
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) repos.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) CreateUser(ctx context.Context, user *entity.User) (int, error) {
	query := sq.Insert("users").
		Columns("name", "email", "password", "created_at").
		Values(user.Name, user.Email, user.Password, user.CreatedAt).
		RunWith(r.db).
		PlaceholderFormat(sq.Question)

	res, err := query.ExecContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

func (r *authRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := sq.Select("id", "username", "email", "password").
		From("users").
		Where(sq.Eq{"email": email}).
		RunWith(r.db).
		PlaceholderFormat(sq.Question)

	row := query.QueryRowContext(ctx)
	user := &entity.User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return user, nil
}
