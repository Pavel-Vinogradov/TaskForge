package usecase

import (
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/domain/repos"
	"TaskForge/internal/interfaces/auth"
	"TaskForge/pkg/jwt"
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

func (u *AuthUseCase) Register(ctx context.Context, req auth.RegisterRequest) (*auth.ResponseAuth, error) {
	existing, _ := u.repo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Name:      req.Username,
		Email:     req.Email,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}

	id, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := jwt.GenerateJWT(id)
	if err != nil {
		return nil, err
	}

	return &auth.ResponseAuth{
		UserID: id,
		Token:  token,
	}, nil
}

func (u *AuthUseCase) Login(ctx context.Context, req auth.LoginRequest) (*auth.ResponseAuth, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := jwt.GenerateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &auth.ResponseAuth{
		UserID: user.ID,
		Token:  token,
	}, nil
}
