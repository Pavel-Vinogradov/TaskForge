package handler

import (
	"TaskForge/internal/interfaces/common"
	"TaskForge/internal/interfaces/task"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	usecase task.UseCaseTask
}

func NewTaskHandler(u task.UseCaseTask) *TaskHandler {
	return &TaskHandler{usecase: u}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req task.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.usecase.CreateTask(c.Request.Context(), req)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *TaskHandler) ListTask(c *gin.Context) {
	res, err := h.usecase.ListTask(c.Request.Context())
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, res)

}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	res, err := h.usecase.UpdateTask(c.Request.Context())
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, res)

}

func (h *TaskHandler) HistoryTask(c *gin.Context) {
	res, err := h.usecase.HistoryTask(c.Request.Context())
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, res)

}
