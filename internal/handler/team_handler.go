package handler

import (
	"TaskForge/internal/interfaces/team"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	usecase team.UseCaseTeam
}

func NewTeamHandler(u team.UseCaseTeam) *TeamHandler {
	return &TeamHandler{usecase: u}
}

// CreateTeam создает новую команду
// @Summary Создать команду
// @Description Создает новую команду и делает пользователя владельцем
// @Tags teams
// @Accept json
// @Produce json
// @Param request body team.CreateTeamRequest true "Данные для создания команды"
// @Success 201 {object} team.CreateTeamResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/teams [post]
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "user_id", userID)

	var req team.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.usecase.CreateTeam(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// ListTeams получает список команд пользователя
// @Summary Получить команды пользователя
// @Description Возвращает список команд, в которых состоит пользователь
// @Tags teams
// @Produce json
// @Success 200 {object} team.ListTeamsResponse
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/teams [get]
func (h *TeamHandler) ListTeams(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Set user_id in request context for usecase
	ctx := context.WithValue(c.Request.Context(), "user_id", userID)

	res, err := h.usecase.ListTeams(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// InviteUser приглашает пользователя в команду
// @Summary Пригласить в команду
// @Description Приглашает пользователя в команду (только для owner/admin)
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "ID команды"
// @Param request body team.InviteUserRequest true "Данные для приглашения"
// @Success 201 {object} team.InviteUserResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/teams/{id}/invite [post]
func (h *TeamHandler) InviteUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	ctx := context.WithValue(c.Request.Context(), "user_id", userID)

	teamIDStr := c.Param("id")
	teamID := 0
	if _, err := fmt.Sscanf(teamIDStr, "%d", &teamID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}

	var req team.InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.usecase.InviteUser(ctx, teamID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}
