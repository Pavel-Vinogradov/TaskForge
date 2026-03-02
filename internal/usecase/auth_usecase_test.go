package usecase

import (
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/interfaces/auth"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthRepo struct {
	mock.Mock
}

func (m *mockAuthRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *mockAuthRepo) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	repo := new(mockAuthRepo)
	uc := NewAuthUseCase(repo)

	ctx := context.Background()
	req := auth.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	repo.On("GetUserByEmail", ctx, "test@example.com").Return(nil, nil)
	repo.On("CreateUser", ctx, mock.AnythingOfType("*entity.User")).Return(&entity.User{
		ID:        1,
		Name:      "testuser",
		Email:     "test@example.com",
		Password:  "$2a$10$hashedpassword",
		CreatedAt: time.Now(),
	}, nil)

	user, err := uc.Register(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "testuser", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.NotEmpty(t, user.Password)

	repo.AssertExpectations(t)
}

func TestRegister_UserExists(t *testing.T) {
	repo := new(mockAuthRepo)
	uc := NewAuthUseCase(repo)

	ctx := context.Background()
	req := auth.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	existingUser := &entity.User{
		ID:    1,
		Name:  "existinguser",
		Email: "test@example.com",
	}

	repo.On("GetUserByEmail", ctx, "test@example.com").Return(existingUser, nil)

	user, err := uc.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user already exists")
	repo.AssertExpectations(t)
}

func TestRegister_GetUserError(t *testing.T) {
	repo := new(mockAuthRepo)
	uc := NewAuthUseCase(repo)

	ctx := context.Background()
	req := auth.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	repo.On("GetUserByEmail", ctx, "test@example.com").Return(nil, errors.New("database error"))

	user, err := uc.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "database error")
	repo.AssertExpectations(t)
}

func TestRegister_CreateUserError(t *testing.T) {
	repo := new(mockAuthRepo)
	uc := NewAuthUseCase(repo)

	ctx := context.Background()
	req := auth.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	repo.On("GetUserByEmail", ctx, "test@example.com").Return(nil, nil)
	repo.On("CreateUser", ctx, mock.AnythingOfType("*entity.User")).Return(nil, errors.New("failed to create user"))

	user, err := uc.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to create user")
	repo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	repo := new(mockAuthRepo)
	uc := NewAuthUseCase(repo)

	ctx := context.Background()
	req := auth.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	existingUser := &entity.User{
		ID:       1,
		Name:     "testuser",
		Email:    "test@example.com",
		Password: string(hash),
	}

	repo.On("GetUserByEmail", ctx, "test@example.com").Return(existingUser, nil)

	userID, err := uc.Login(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, userID)
	repo.AssertExpectations(t)
}

func TestLogin_EmptyCredentials(t *testing.T) {
	repo := new(mockAuthRepo)
	uc := NewAuthUseCase(repo)

	ctx := context.Background()

	testCases := []auth.LoginRequest{
		{Email: "", Password: "password123"},
		{Email: "test@example.com", Password: ""},
		{Email: "", Password: ""},
	}

	for _, req := range testCases {
		userID, err := uc.Login(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, 0, userID)
		assert.Contains(t, err.Error(), "email and password are required")
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := new(mockAuthRepo)
	uc := NewAuthUseCase(repo)

	ctx := context.Background()
	req := auth.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	repo.On("GetUserByEmail", ctx, "nonexistent@example.com").Return(nil, errors.New("user not found"))

	userID, err := uc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.Contains(t, err.Error(), "invalid credentials")
	repo.AssertExpectations(t)
}

func TestLogin_UserNil(t *testing.T) {
	repo := new(mockAuthRepo)
	uc := NewAuthUseCase(repo)

	ctx := context.Background()
	req := auth.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	repo.On("GetUserByEmail", ctx, "nonexistent@example.com").Return(nil, nil)

	userID, err := uc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.Contains(t, err.Error(), "invalid credentials")
	repo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	repo := new(mockAuthRepo)
	uc := NewAuthUseCase(repo)

	ctx := context.Background()
	req := auth.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	existingUser := &entity.User{
		ID:       1,
		Name:     "testuser",
		Email:    "test@example.com",
		Password: string(hash),
	}

	repo.On("GetUserByEmail", ctx, "test@example.com").Return(existingUser, nil)

	userID, err := uc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.Contains(t, err.Error(), "invalid credentials")
	repo.AssertExpectations(t)
}

func TestLogin_GetUserError(t *testing.T) {
	repo := new(mockAuthRepo)
	uc := NewAuthUseCase(repo)

	ctx := context.Background()
	req := auth.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	repo.On("GetUserByEmail", ctx, "test@example.com").Return(nil, errors.New("database error"))

	userID, err := uc.Login(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.Contains(t, err.Error(), "invalid credentials")
	repo.AssertExpectations(t)
}
