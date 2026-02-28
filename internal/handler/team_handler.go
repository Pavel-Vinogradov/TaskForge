package handler

import (
	"TaskForge/internal/interfaces/common"
	"TaskForge/internal/interfaces/team"
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
// @Success 201 {object} common.Response{data=team.CreateTeamResponse}
// @Failure 400 {object} common.Response
// @Failure 500 {object} common.Response
// @Security BearerAuth
// @Router /api/v1/teams [post]
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req team.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.usecase.CreateTeam(c.Request.Context(), req)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, common.Response{
		Success: true,
		Data:    res,
	})
}

// ListTeams получает список команд пользователя
// @Summary Получить команды пользователя
// @Description Возвращает список команд, в которых состоит пользователь
// @Tags teams
// @Produce json
// @Success 200 {object} common.Response{data=[]team.TeamInfo}
// @Failure 500 {object} common.Response
// @Security BearerAuth
// @Router /api/v1/teams [get]
func (h *TeamHandler) ListTeams(c *gin.Context) {
	res, err := h.usecase.ListTeams(c.Request.Context())
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, common.Response{
		Success: true,
		Data:    res,
	})
}

// InviteUser приглашает пользователя в команду
// @Summary Пригласить в команду
// @Description Приглашает пользователя в команду (только для owner/admin)
// @Tags teams
// @Accept json
// @Produce json
// @Param id path int true "ID команды"
// @Param request body team.InviteUserRequest true "Данные для приглашения"
// @Success 201 {object} common.Response{data=team.InviteUserResponse}
// @Failure 400 {object} common.Response
// @Failure 403 {object} common.Response
// @Failure 500 {object} common.Response
// @Security BearerAuth
// @Router /api/v1/teams/{id}/invite [post]
func (h *TeamHandler) InviteUser(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID := 0
	if _, err := fmt.Sscanf(teamIDStr, "%d", &teamID); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	var req team.InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.usecase.InviteUser(c.Request.Context(), teamID, req)
	if err != nil {
		common.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, common.Response{
		Success: true,
		Data:    res,
	})
}
