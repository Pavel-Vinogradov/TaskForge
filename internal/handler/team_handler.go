package handler

import (
	"TaskForge/internal/interfaces/team"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	usecase team.UseCaseTeam
}

func (h TeamHandler) CreateHandler(context *gin.Context) {

}

func (h TeamHandler) ListHandler(context *gin.Context) {

}

func (h TeamHandler) InviteHandler(context *gin.Context) {

}

func NewTeamHandler(u team.UseCaseTeam) *TeamHandler {
	return &TeamHandler{usecase: u}
}
