package handler

import (
	"TaskForge/internal/interfaces/task"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	usecase task.UseCaseTask
}

func NewTaskHandler(u task.UseCaseTask) *TaskHandler {
	return &TaskHandler{usecase: u}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {

}

func (h *TaskHandler) ListTask(c *gin.Context) {

}

func (h *TaskHandler) UpdateTask(c *gin.Context) {

}

func (h *TaskHandler) HistoryTask(c *gin.Context) {

}
