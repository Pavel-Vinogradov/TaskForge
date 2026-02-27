package team

import "time"

type CreateTeamRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100"`
}

type CreateTeamResponse struct {
	ID        int
	Name      string
	Role      string
	CreatedBy int
	CreatedAt time.Time
}

type UseCaseTeam interface {
}
