package repos

import (
	"TaskForge/internal/domain/entity"
	"context"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}
