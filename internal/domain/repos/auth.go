package repos

import (
	"TaskForge/internal/domain/entity"
	"context"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}
