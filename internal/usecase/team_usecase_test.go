package usecase

import (
	"TaskForge/internal/contextkeys"
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/interfaces/team"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTeamRepo struct {
	mock.Mock
}

func (m *mockTeamRepo) CreateTeam(ctx context.Context, t entity.Team) (entity.Team, error) {
	args := m.Called(ctx, t)
	return args.Get(0).(entity.Team), args.Error(1)
}

func (m *mockTeamRepo) AddTeamMember(ctx context.Context, member entity.TeamMember) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *mockTeamRepo) GetUserTeams(ctx context.Context, userID int) ([]team.TeamWithMembership, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]team.TeamWithMembership), args.Error(1)
}

func (m *mockTeamRepo) GetTeamMember(ctx context.Context, teamID, userID int) (entity.TeamMember, error) {
	args := m.Called(ctx, teamID, userID)
	return args.Get(0).(entity.TeamMember), args.Error(1)
}

func TestCreateTeam_Success(t *testing.T) {
	repo := new(mockTeamRepo)
	uc := NewTeamUseCase(repo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	req := team.CreateTeamRequest{
		Name: "Test Team",
	}

	expectedTeam := entity.Team{
		ID:        1,
		Name:      "Test Team",
		CreatedBy: 42,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("CreateTeam", ctx, mock.AnythingOfType("entity.Team")).Return(expectedTeam, nil)
	repo.On("AddTeamMember", ctx, mock.AnythingOfType("entity.TeamMember")).Return(nil)

	res, err := uc.CreateTeam(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.ID)
	assert.Equal(t, "Test Team", res.Name)
	assert.Equal(t, "owner", res.Role)
	assert.Equal(t, 42, res.CreatedBy)
	repo.AssertExpectations(t)
}

func TestCreateTeam_NoUser(t *testing.T) {
	repo := new(mockTeamRepo)
	uc := NewTeamUseCase(repo)

	ctx := context.Background()

	req := team.CreateTeamRequest{
		Name: "Test Team",
	}

	_, err := uc.CreateTeam(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not authenticated")
}

func TestCreateTeam_CreateTeamError(t *testing.T) {
	repo := new(mockTeamRepo)
	uc := NewTeamUseCase(repo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	req := team.CreateTeamRequest{
		Name: "Test Team",
	}

	repo.On("CreateTeam", ctx, mock.AnythingOfType("entity.Team")).Return(entity.Team{}, errors.New("database error"))

	_, err := uc.CreateTeam(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	repo.AssertExpectations(t)
}

func TestCreateTeam_AddMemberError(t *testing.T) {
	repo := new(mockTeamRepo)
	uc := NewTeamUseCase(repo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	req := team.CreateTeamRequest{
		Name: "Test Team",
	}

	expectedTeam := entity.Team{
		ID:        1,
		Name:      "Test Team",
		CreatedBy: 42,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.On("CreateTeam", ctx, mock.AnythingOfType("entity.Team")).Return(expectedTeam, nil)
	repo.On("AddTeamMember", ctx, mock.AnythingOfType("entity.TeamMember")).Return(errors.New("failed to add member"))

	_, err := uc.CreateTeam(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to add member")
	repo.AssertExpectations(t)
}

func TestListTeams_Success(t *testing.T) {
	repo := new(mockTeamRepo)
	uc := NewTeamUseCase(repo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	teams := []team.TeamWithMembership{
		{
			Team: entity.Team{
				ID:        1,
				Name:      "Team 1",
				CreatedAt: time.Now(),
			},
			Role:     entity.RoleOwner,
			JoinedAt: time.Now(),
		},
		{
			Team: entity.Team{
				ID:        2,
				Name:      "Team 2",
				CreatedAt: time.Now(),
			},
			Role:     entity.RoleMember,
			JoinedAt: time.Now(),
		},
	}

	repo.On("GetUserTeams", ctx, 42).Return(teams, nil)

	res, err := uc.ListTeams(ctx)

	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.Equal(t, 1, res[0].ID)
	assert.Equal(t, "Team 1", res[0].Name)
	assert.Equal(t, "owner", res[0].Role)
	assert.Equal(t, 2, res[1].ID)
	assert.Equal(t, "Team 2", res[1].Name)
	assert.Equal(t, "member", res[1].Role)
	repo.AssertExpectations(t)
}

func TestListTeams_NoUser(t *testing.T) {
	repo := new(mockTeamRepo)
	uc := NewTeamUseCase(repo)

	ctx := context.Background()

	_, err := uc.ListTeams(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not authenticated")
}

func TestListTeams_RepoError(t *testing.T) {
	repo := new(mockTeamRepo)
	uc := NewTeamUseCase(repo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	repo.On("GetUserTeams", ctx, 42).Return([]team.TeamWithMembership{}, errors.New("database error"))

	_, err := uc.ListTeams(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	repo.AssertExpectations(t)
}

func TestInviteUser_Success(t *testing.T) {
	repo := new(mockTeamRepo)
	uc := NewTeamUseCase(repo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	req := team.InviteUserRequest{
		UserID: 123,
	}

	inviterMember := entity.TeamMember{
		TeamID: 1,
		UserID: 42,
		Role:   entity.RoleOwner,
	}

	repo.On("GetTeamMember", ctx, 1, 42).Return(inviterMember, nil)
	repo.On("GetTeamMember", ctx, 1, 123).Return(entity.TeamMember{}, errors.New("not found"))
	repo.On("AddTeamMember", ctx, mock.AnythingOfType("entity.TeamMember")).Return(nil)

	res, err := uc.InviteUser(ctx, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.TeamID)
	assert.Equal(t, 123, res.UserID)
	assert.Equal(t, "member", res.Role)
	repo.AssertExpectations(t)
}

func TestInviteUser_NoUser(t *testing.T) {
	repo := new(mockTeamRepo)
	uc := NewTeamUseCase(repo)

	ctx := context.Background()

	req := team.InviteUserRequest{
		UserID: 123,
	}

	_, err := uc.InviteUser(ctx, 1, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not authenticated")
}

func TestInviteUser_InsufficientPermissions(t *testing.T) {
	repo := new(mockTeamRepo)
	uc := NewTeamUseCase(repo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	req := team.InviteUserRequest{
		UserID: 123,
	}

	inviterMember := entity.TeamMember{
		TeamID: 1,
		UserID: 42,
		Role:   entity.RoleMember,
	}

	repo.On("GetTeamMember", ctx, 1, 42).Return(inviterMember, nil)

	_, err := uc.InviteUser(ctx, 1, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient permissions")
	repo.AssertExpectations(t)
}

func TestInviteUser_UserAlreadyMember(t *testing.T) {
	repo := new(mockTeamRepo)
	uc := NewTeamUseCase(repo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	req := team.InviteUserRequest{
		UserID: 123,
	}

	inviterMember := entity.TeamMember{
		TeamID: 1,
		UserID: 42,
		Role:   entity.RoleOwner,
	}

	existingMember := entity.TeamMember{
		TeamID: 1,
		UserID: 123,
		Role:   entity.RoleMember,
	}

	repo.On("GetTeamMember", ctx, 1, 42).Return(inviterMember, nil)
	repo.On("GetTeamMember", ctx, 1, 123).Return(existingMember, nil)

	_, err := uc.InviteUser(ctx, 1, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user is already a team member")
	repo.AssertExpectations(t)
}

func TestInviteUser_AddMemberError(t *testing.T) {
	repo := new(mockTeamRepo)
	uc := NewTeamUseCase(repo)

	ctx := context.WithValue(context.Background(), contextkeys.UserIDKey, 42)

	req := team.InviteUserRequest{
		UserID: 123,
	}

	inviterMember := entity.TeamMember{
		TeamID: 1,
		UserID: 42,
		Role:   entity.RoleOwner,
	}

	repo.On("GetTeamMember", ctx, 1, 42).Return(inviterMember, nil)
	repo.On("GetTeamMember", ctx, 1, 123).Return(entity.TeamMember{}, errors.New("not found"))
	repo.On("AddTeamMember", ctx, mock.AnythingOfType("entity.TeamMember")).Return(errors.New("failed to add member"))

	_, err := uc.InviteUser(ctx, 1, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to add member")
	repo.AssertExpectations(t)
}
