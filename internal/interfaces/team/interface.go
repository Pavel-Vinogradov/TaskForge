package team

import (
	"TaskForge/internal/domain/entity"
	"context"
	"time"
)

type CreateTeamRequest struct {
	Name string `json:"name" binding:"required,min=3,max=100" example:"My Team"`
}

type CreateTeamResponse struct {
	ID        int       `json:"id" example:"1"`
	Name      string    `json:"name" example:"My Team"`
	Role      string    `json:"role" example:"owner"`
	CreatedBy int       `json:"created_by" example:"123"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"`
}

type ListTeamsResponse struct {
	Teams []TeamInfo `json:"teams"`
}

type TeamInfo struct {
	ID        int       `json:"id" example:"1"`
	Name      string    `json:"name" example:"My Team"`
	Role      string    `json:"role" example:"member"`
	JoinedAt  time.Time `json:"joined_at" example:"2024-01-01T12:00:00Z"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

type InviteUserRequest struct {
	UserID int `json:"user_id" binding:"required" example:"456"`
}

type InviteUserResponse struct {
	TeamID   int       `json:"team_id" example:"1"`
	UserID   int       `json:"user_id" example:"456"`
	Role     string    `json:"role" example:"member"`
	JoinedAt time.Time `json:"joined_at" example:"2024-01-01T12:00:00Z"`
}

type UseCaseTeam interface {
	CreateTeam(ctx context.Context, req CreateTeamRequest) (CreateTeamResponse, error)
	ListTeams(ctx context.Context) (ListTeamsResponse, error)
	InviteUser(ctx context.Context, teamID int, req InviteUserRequest) (InviteUserResponse, error)
}

type RepositoryTeam interface {
	CreateTeam(ctx context.Context, team entity.Team) (entity.Team, error)
	GetUserTeams(ctx context.Context, userID int) ([]TeamWithMembership, error)
	GetTeamMember(ctx context.Context, teamID, userID int) (entity.TeamMember, error)
	AddTeamMember(ctx context.Context, member entity.TeamMember) error
}

type TeamWithMembership struct {
	Team     entity.Team
	Member   entity.TeamMember
	Role     entity.Role
	JoinedAt time.Time
}
