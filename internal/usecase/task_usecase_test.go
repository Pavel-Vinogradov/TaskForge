package usecase

import (
	"TaskForge/internal/contextkeys"
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/domain/repos"
	"TaskForge/internal/interfaces/task"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

type mockTaskRepo struct{}

func (m *mockTaskRepo) CreateTask(ctx context.Context, t entity.Task) (entity.Task, error) {
	t.Id = 1
	return t, nil
}

func (m *mockTaskRepo) ListTasks(ctx context.Context, f task.TaskFilters) ([]entity.Task, int64, error) {
	return nil, 0, nil
}
func (m *mockTaskRepo) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	if id == 999 {
		return entity.Task{}, errors.New("task not found")
	}
	return entity.Task{
		Id:          id,
		Title:       "Existing task",
		Description: "Description",
		Status:      entity.StatusTodo,
		TeamID:      10,
		CreatedBy:   42,
		AssigneeID:  5,
		CreatedAt:   time.Now(),
	}, nil
}
func (m *mockTaskRepo) UpdateTask(ctx context.Context, t entity.Task) (entity.Task, error) {
	return t, nil
}

type mockErrorTaskRepo struct{}

func (m *mockErrorTaskRepo) CreateTask(ctx context.Context, t entity.Task) (entity.Task, error) {
	return entity.Task{}, errors.New("database error")
}

func (m *mockErrorTaskRepo) ListTasks(ctx context.Context, f task.TaskFilters) ([]entity.Task, int64, error) {
	return nil, 0, nil
}
func (m *mockErrorTaskRepo) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	return entity.Task{}, nil
}
func (m *mockErrorTaskRepo) UpdateTask(ctx context.Context, t entity.Task) (entity.Task, error) {
	return entity.Task{}, nil
}

func TestCreateTask_Success(t *testing.T) {
	repo := &mockTaskRepo{}
	var historyRepo repos.TaskHistoryRepository = nil
	var redisClient *redis.Client = nil

	uc := NewTaskUseCase(repo, historyRepo, redisClient)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	req := task.CreateTaskRequest{
		Title:       "Test task",
		Description: "Test description",
		TeamID:      10,
		AssigneeID:  5,
	}

	res, err := uc.CreateTask(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.ID)
	assert.Equal(t, "Test task", res.Title)
	assert.Equal(t, "Test description", res.Description)
	assert.Equal(t, "todo", res.Status)
	assert.Equal(t, 10, res.TeamID)
	assert.Equal(t, 42, res.CreatedBy)
	assert.Equal(t, 5, res.AssigneeID)
}

func TestCreateTask_NoUser(t *testing.T) {
	repo := &mockTaskRepo{}
	uc := NewTaskUseCase(repo, nil, nil)

	ctx := context.Background()

	req := task.CreateTaskRequest{}

	_, err := uc.CreateTask(ctx, req)

	assert.Error(t, err)
}

func TestCreateTask_RepoError(t *testing.T) {
	repo := &mockErrorTaskRepo{}
	uc := NewTaskUseCase(repo, nil, nil)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	req := task.CreateTaskRequest{
		Title:       "Test task",
		Description: "Test description",
		TeamID:      10,
		AssigneeID:  5,
	}

	_, err := uc.CreateTask(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestUpdateTask_Success(t *testing.T) {
	repo := &mockTaskRepo{}
	uc := NewTaskUseCase(repo, nil, nil)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	newTitle := "Updated title"
	req := task.UpdateTaskRequest{
		Title: &newTitle,
	}

	res, err := uc.UpdateTask(ctx, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.ID)
	assert.Equal(t, "Updated title", res.Title)
}

func TestUpdateTask_NotFound(t *testing.T) {
	repo := &mockTaskRepo{}
	uc := NewTaskUseCase(repo, nil, nil)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	req := task.UpdateTaskRequest{}

	_, err := uc.UpdateTask(ctx, 999, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task not found")
}

func TestUpdateTask_NoPermission(t *testing.T) {
	repo := &mockTaskRepo{}
	uc := NewTaskUseCase(repo, nil, nil)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 99) // Different user

	req := task.UpdateTaskRequest{}

	_, err := uc.UpdateTask(ctx, 1, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient permissions")
}
