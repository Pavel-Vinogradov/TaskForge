package usecase

import (
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/domain/repos"
	"TaskForge/internal/interfaces/auth"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	repo repos.AuthRepository
}

func NewAuthUseCase(repo repos.AuthRepository) *AuthUseCase {
	return &AuthUseCase{
		repo: repo,
	}
}

func (u *AuthUseCase) Register(ctx context.Context, req auth.RegisterRequest) (int, error) {
	existing, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := &entity.User{
		Name:      req.Username,
		Email:     req.Email,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}

	id, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *AuthUseCase) Login(ctx context.Context, req auth.LoginRequest) (int, error) {
	if req.Email == "" || req.Password == "" {
		return 0, errors.New("email and password are required")
	}

	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return 0, errors.New("invalid credentials")
	}

	if user == nil {
		return 0, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return 0, errors.New("invalid credentials")
	}

	return user.ID, nil
}
