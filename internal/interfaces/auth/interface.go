package auth

import "context"

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type ResponseAuth struct {
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
}

type UseCaseAuth interface {
	Register(ctx context.Context, req RegisterRequest) (int, error)
	Login(ctx context.Context, req LoginRequest) (int, error)
}
