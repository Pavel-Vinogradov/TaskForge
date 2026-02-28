package repos

import (
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/interfaces/team"
	"context"
)

type TeamRepository interface {
	CreateTeam(ctx context.Context, team entity.Team) (entity.Team, error)
	GetUserTeams(ctx context.Context, userID int) ([]team.TeamWithMembership, error)
	GetTeamMember(ctx context.Context, teamID, userID int) (entity.TeamMember, error)
	AddTeamMember(ctx context.Context, member entity.TeamMember) error
}
