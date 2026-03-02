package usecase

import (
	"TaskForge/internal/contextkeys"
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/domain/repos"
	"TaskForge/internal/interfaces/team"
	"context"
	"errors"
	"time"
)

type TeamUseCase struct {
	repo repos.TeamRepository
}

func NewTeamUseCase(repo repos.TeamRepository) *TeamUseCase {
	return &TeamUseCase{repo: repo}
}

func (uc *TeamUseCase) CreateTeam(ctx context.Context, req team.CreateTeamRequest) (team.CreateTeamResponse, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey).(int)
	if !ok {
		return team.CreateTeamResponse{}, errors.New("user not authenticated")
	}

	teamEntity := entity.Team{
		Name:      req.Name,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdTeam, err := uc.repo.CreateTeam(ctx, teamEntity)
	if err != nil {
		return team.CreateTeamResponse{}, err
	}

	teamMember := entity.TeamMember{
		TeamID:   createdTeam.ID,
		UserID:   userID,
		Role:     entity.RoleOwner,
		JoinedAt: time.Now(),
	}

	err = uc.repo.AddTeamMember(ctx, teamMember)
	if err != nil {
		return team.CreateTeamResponse{}, err
	}

	return team.CreateTeamResponse{
		ID:        createdTeam.ID,
		Name:      createdTeam.Name,
		Role:      string(entity.RoleOwner),
		CreatedBy: createdTeam.CreatedBy,
		CreatedAt: createdTeam.CreatedAt,
	}, nil
}

func (uc *TeamUseCase) ListTeams(ctx context.Context) ([]team.TeamInfo, error) {
	userID, ok := ctx.Value(contextkeys.UserIDKey).(int)
	if !ok {
		return nil, errors.New("user not authenticated")
	}

	teams, err := uc.repo.GetUserTeams(ctx, userID)
	if err != nil {
		return nil, err
	}

	var teamInfos []team.TeamInfo
	for _, t := range teams {
		teamInfos = append(teamInfos, team.TeamInfo{
			ID:        t.Team.ID,
			Name:      t.Team.Name,
			Role:      string(t.Role),
			JoinedAt:  t.JoinedAt,
			CreatedAt: t.Team.CreatedAt,
		})
	}

	return teamInfos, nil
}

func (uc *TeamUseCase) InviteUser(ctx context.Context, teamID int, req team.InviteUserRequest) (team.InviteUserResponse, error) {
	inviterID, ok := ctx.Value(contextkeys.UserIDKey).(int)
	if !ok {
		return team.InviteUserResponse{}, errors.New("user not authenticated")
	}

	inviterMember, err := uc.repo.GetTeamMember(ctx, teamID, inviterID)
	if err != nil {
		return team.InviteUserResponse{}, errors.New("failed to check team membership")
	}

	if inviterMember.Role != entity.RoleOwner && inviterMember.Role != entity.RoleAdmin {
		return team.InviteUserResponse{}, errors.New("insufficient permissions to invite users")
	}

	_, err = uc.repo.GetTeamMember(ctx, teamID, req.UserID)
	if err == nil {
		return team.InviteUserResponse{}, errors.New("user is already a team member")
	}

	newMember := entity.TeamMember{
		TeamID:   teamID,
		UserID:   req.UserID,
		Role:     entity.RoleMember,
		JoinedAt: time.Now(),
	}

	err = uc.repo.AddTeamMember(ctx, newMember)
	if err != nil {
		return team.InviteUserResponse{}, err
	}

	return team.InviteUserResponse{
		TeamID:   teamID,
		UserID:   req.UserID,
		Role:     string(entity.RoleMember),
		JoinedAt: newMember.JoinedAt,
	}, nil
}
