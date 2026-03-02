package handler

import (
	"TaskForge/internal/interfaces/common"
	"TaskForge/internal/interfaces/task"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	usecase task.UseCaseTask
}

func NewTaskHandler(u task.UseCaseTask) *TaskHandler {
	return &TaskHandler{usecase: u}
}

// CreateTask godoc
// @Summary Create a new task
// @Description Creates a new task with the provided details. Requires authentication.
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body task.CreateTaskRequest true "Task creation request"
// @Success 201 {object} common.Response{Data=task.ResponseTask} "Task created successfully"
// @Failure 400 {object} common.Response "Invalid request body"
// @Failure 500 {object} common.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/tasks [post]
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
	c.JSON(http.StatusCreated, common.Response{
		Success: true,
		Data:    res,
	})
}

// ListTask godoc
// @Summary List tasks
// @Description Retrieves a paginated list of tasks with optional filtering. Requires authentication.
// @Tags tasks
// @Accept json
// @Produce json
// @Param team_id query int false "Filter by team ID"
// @Param status query string false "Filter by status (todo, in_progress, done)"
// @Param assignee_id query int false "Filter by assignee ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} common.PaginationResponse{Data=[]entity.Task} "Tasks retrieved successfully"
// @Failure 400 {object} common.Response "Invalid query parameters"
// @Failure 500 {object} common.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/tasks [get]
func (h *TaskHandler) ListTask(c *gin.Context) {
	var req task.TaskListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	res, err := h.usecase.ListTask(c.Request.Context(), req)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, common.PaginationResponse{
		Success: true,
		Data:    res.Tasks,
		Page:    req.Page,
		Limit:   req.Limit,
		Total:   int(res.Total),
	})

}

// UpdateTask godoc
// @Summary Update a task
// @Description Updates an existing task. Requires authentication.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task Id"
// @Param request body task.UpdateTaskRequest false "Task update request"
// @Success 200 {object} common.Response{data=task.ResponseTask}
// @Failure 400 {object} common.Response "Invalid request body"
// @Failure 403 {object} common.Response "Insufficient permissions"
// @Failure 404 {object} common.Response "Task not found"
// @Failure 500 {object} common.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid task id"))
		return
	}

	var req task.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.usecase.UpdateTask(c.Request.Context(), taskID, req)
	if err != nil {
		if errors.Is(err, errors.New("task not found")) {
			common.ErrorResponse(c, http.StatusNotFound, err)
			return
		}
		if errors.Is(err, errors.New("insufficient permissions")) {
			common.ErrorResponse(c, http.StatusForbidden, err)
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, common.Response{Success: true, Data: res})

}

// HistoryTask godoc
// @Summary Get task history
// @Description Retrieves the history of changes for a specific task. Requires authentication.
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task Id"
// @Success 200 {object} common.Response{data=[]entity.TaskHistory}
// @Failure 400 {object} common.Response "Invalid task ID"
// @Failure 403 {object} common.Response "Insufficient permissions"
// @Failure 404 {object} common.Response "Task not found"
// @Failure 500 {object} common.Response "Internal server error"
// @Security ApiKeyAuth
// @Router /api/v1/tasks/{id}/history [get]
func (h *TaskHandler) HistoryTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, errors.New("invalid task id"))
		return
	}

	res, err := h.usecase.HistoryTask(c.Request.Context(), taskID)
	if err != nil {
		if errors.Is(err, errors.New("task not found")) {
			common.ErrorResponse(c, http.StatusNotFound, err)
			return
		}
		if errors.Is(err, errors.New("insufficient permissions")) {
			common.ErrorResponse(c, http.StatusForbidden, err)
			return
		}
		common.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, common.Response{Success: true, Data: res})

}
