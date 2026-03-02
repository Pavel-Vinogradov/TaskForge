package usecase

import (
	"TaskForge/internal/contextkeys"
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/domain/repos"
	"TaskForge/internal/interfaces/task"
	"TaskForge/internal/interfaces/task_history"
	"context"
	"errors"
	"time"
)

type TaskUseCase struct {
	repo        repos.TaskRepository
	historyRepo repos.TaskHistoryRepository
	observers   []task_history.TaskObserver
}

func NewTaskUseCase(repo repos.TaskRepository, historyRepo repos.TaskHistoryRepository) *TaskUseCase {
	return &TaskUseCase{
		repo:        repo,
		historyRepo: historyRepo,
		observers:   []task_history.TaskObserver{},
	}
}

func (uc *TaskUseCase) AddObserver(observer task_history.TaskObserver) {
	uc.observers = append(uc.observers, observer)
}

func (uc *TaskUseCase) CreateTask(ctx context.Context, req task.CreateTaskRequest) (task.ResponseTask, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey).(int)
	if !ok {
		return task.ResponseTask{}, errors.New("user not authenticated")
	}

	taskEntity := entity.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      entity.StatusTodo,
		TeamID:      req.TeamID,
		CreatedBy:   userID,
		AssigneeID:  req.AssigneeID,
		CreatedAt:   time.Now(),
	}

	createdTask, err := uc.repo.CreateTask(ctx, taskEntity)
	if err != nil {
		return task.ResponseTask{}, err
	}

	for _, observer := range uc.observers {
		observer.OnTaskCreated(ctx, createdTask, userID)
	}

	return task.ResponseTask{
		ID:          createdTask.Id,
		Title:       createdTask.Title,
		Description: createdTask.Description,
		Status:      string(createdTask.Status),
		TeamID:      createdTask.TeamID,
		CreatedBy:   createdTask.CreatedBy,
		AssigneeID:  createdTask.AssigneeID,
		CreatedAt:   createdTask.CreatedAt,
	}, nil
}

func (uc *TaskUseCase) ListTask(ctx context.Context, req task.TaskListRequest) (task.TaskListResult, error) {
	_, ok := ctx.Value(contextkeys.UserIDKey).(int)
	if !ok {
		return task.TaskListResult{}, errors.New("user not authenticated")
	}

	filters := task.TaskFilters{
		TeamID:     req.TeamID,
		Status:     req.Status,
		AssigneeID: req.AssigneeID,
		Page:       req.Page,
		Limit:      req.Limit,
	}

	tasks, total, err := uc.repo.ListTasks(ctx, filters)
	if err != nil {
		return task.TaskListResult{}, err
	}
	return task.TaskListResult{
		Tasks: tasks,
		Total: total,
	}, nil
}

func (uc *TaskUseCase) UpdateTask(ctx context.Context, taskID int, req task.UpdateTaskRequest) (task.ResponseTask, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey).(int)
	if !ok {
		return task.ResponseTask{}, errors.New("user not authenticated")
	}

	existingTask, err := uc.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return task.ResponseTask{}, errors.New("task not found")
	}

	if existingTask.CreatedBy != userID && existingTask.AssigneeID != userID {
		return task.ResponseTask{}, errors.New("insufficient permissions")
	}

	originalTask := existingTask

	if req.Title != nil {
		existingTask.Title = *req.Title
	}
	if req.Description != nil {
		existingTask.Description = *req.Description
	}
	if req.Status != nil {
		existingTask.Status = entity.TaskStatus(*req.Status)
	}
	if req.AssigneeID != nil {
		existingTask.AssigneeID = *req.AssigneeID
	}

	existingTask.UpdatedAt = time.Now()

	updatedTask, err := uc.repo.UpdateTask(ctx, existingTask)
	if err != nil {
		return task.ResponseTask{}, err
	}

	for _, observer := range uc.observers {
		observer.OnTaskUpdated(ctx, originalTask, updatedTask, userID)
	}

	return task.ResponseTask{
		ID:          updatedTask.Id,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		Status:      string(updatedTask.Status),
		TeamID:      updatedTask.TeamID,
		CreatedBy:   updatedTask.CreatedBy,
		AssigneeID:  updatedTask.AssigneeID,
		CreatedAt:   updatedTask.CreatedAt,
	}, nil
}

func (uc *TaskUseCase) HistoryTask(ctx context.Context, taskID int) ([]entity.TaskHistory, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey).(int)
	if !ok {
		return nil, errors.New("user not authenticated")
	}

	task, err := uc.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, errors.New("task not found")
	}

	if task.CreatedBy != userID && task.AssigneeID != userID {
		return nil, errors.New("insufficient permissions")
	}

	history, err := uc.historyRepo.GetTaskHistory(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return history, nil
}
