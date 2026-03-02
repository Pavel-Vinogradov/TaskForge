package usecase

import (
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/domain/repos"
	"context"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type TaskHistoryObserver struct {
	historyRepo repos.TaskHistoryRepository
}

func NewTaskHistoryObserver(historyRepo repos.TaskHistoryRepository) *TaskHistoryObserver {
	return &TaskHistoryObserver{
		historyRepo: historyRepo,
	}
}

func (h *TaskHistoryObserver) OnTaskCreated(ctx context.Context, task entity.Task, userID int) {
	history := entity.TaskHistory{
		TaskID:    task.Id,
		Field:     "status",
		OldValue:  "",
		NewValue:  string(task.Status),
		ChangedBy: userID,
		ChangedAt: time.Now(),
	}

	if err := h.historyRepo.CreateHistory(ctx, history); err != nil {
		logrus.WithError(err).Error("Failed to create task history for created task")
	}
}

func (h *TaskHistoryObserver) OnTaskUpdated(ctx context.Context, oldTask, newTask entity.Task, userID int) {
	if oldTask.Status != newTask.Status {
		history := entity.TaskHistory{
			TaskID:    newTask.Id,
			Field:     "status",
			OldValue:  string(oldTask.Status),
			NewValue:  string(newTask.Status),
			ChangedBy: userID,
			ChangedAt: time.Now(),
		}

		if err := h.historyRepo.CreateHistory(ctx, history); err != nil {
			logrus.WithError(err).WithField("task_id", newTask.Id).Error("Failed to create task history")
		}
	}

	if oldTask.Title != newTask.Title {
		history := entity.TaskHistory{
			TaskID:    newTask.Id,
			Field:     "title",
			OldValue:  oldTask.Title,
			NewValue:  newTask.Title,
			ChangedBy: userID,
			ChangedAt: time.Now(),
		}

		if err := h.historyRepo.CreateHistory(ctx, history); err != nil {
			logrus.WithError(err).WithField("task_id", newTask.Id).Error("Failed to create task history")
		}
	}

	if oldTask.Description != newTask.Description {
		history := entity.TaskHistory{
			TaskID:    newTask.Id,
			Field:     "description",
			OldValue:  oldTask.Description,
			NewValue:  newTask.Description,
			ChangedBy: userID,
			ChangedAt: time.Now(),
		}

		if err := h.historyRepo.CreateHistory(ctx, history); err != nil {
			logrus.WithError(err).WithField("task_id", newTask.Id).Error("Failed to create task history")
		}
	}

	if oldTask.AssigneeID != newTask.AssigneeID {
		history := entity.TaskHistory{
			TaskID:    newTask.Id,
			Field:     "assignee_id",
			OldValue:  strconv.Itoa(oldTask.AssigneeID),
			NewValue:  strconv.Itoa(newTask.AssigneeID),
			ChangedBy: userID,
			ChangedAt: time.Now(),
		}

		if err := h.historyRepo.CreateHistory(ctx, history); err != nil {
			logrus.WithError(err).WithField("task_id", newTask.Id).Error("Failed to create task history")
		}
	}
}
