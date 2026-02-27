package repository

import (
	"TaskForge/internal/domain/entity"
	"TaskForge/internal/interfaces/team"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type teamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) team.RepositoryTeam {
	return &teamRepository{db: db}
}

func (r *teamRepository) CreateTeam(ctx context.Context, team entity.Team) (entity.Team, error) {
	query := sq.Insert("teams").
		Columns("name", "created_by", "created_at", "updated_at").
		Values(team.Name, team.CreatedBy, team.CreatedAt, team.UpdatedAt).
		RunWith(r.db).
		PlaceholderFormat(sq.Question)

	res, err := query.ExecContext(ctx)
	if err != nil {
		return entity.Team{}, fmt.Errorf("failed to insert team: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return entity.Team{}, fmt.Errorf("failed to get last insert id: %w", err)
	}

	team.ID = int(id)
	return team, nil
}

func (r *teamRepository) GetUserTeams(ctx context.Context, userID int) ([]team.TeamWithMembership, error) {
	query := sq.Select(
		"t.id", "t.name", "t.created_by", "t.created_at", "t.updated_at",
		"tm.user_id", "tm.role", "tm.joined_at",
	).
		From("teams t").
		Join("team_members tm ON t.id = tm.team_id").
		Where(sq.Eq{"tm.user_id": userID}).
		RunWith(r.db).
		PlaceholderFormat(sq.Question)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query user teams: %w", err)
	}
	defer rows.Close()

	var teams []team.TeamWithMembership
	for rows.Next() {
		var t entity.Team
		var tm entity.TeamMember
		var joinedAt time.Time

		err := rows.Scan(
			&t.ID, &t.Name, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
			&tm.UserID, &tm.Role, &joinedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan team row: %w", err)
		}

		tm.TeamID = t.ID
		tm.JoinedAt = joinedAt

		teams = append(teams, team.TeamWithMembership{
			Team:     t,
			Member:   tm,
			Role:     tm.Role,
			JoinedAt: joinedAt,
		})
	}

	return teams, nil
}

func (r *teamRepository) GetTeamMember(ctx context.Context, teamID, userID int) (entity.TeamMember, error) {
	query := sq.Select("team_id", "user_id", "role", "joined_at").
		From("team_members").
		Where(sq.Eq{"team_id": teamID, "user_id": userID}).
		RunWith(r.db).
		PlaceholderFormat(sq.Question)

	row := query.QueryRowContext(ctx)
	var member entity.TeamMember
	err := row.Scan(&member.TeamID, &member.UserID, &member.Role, &member.JoinedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.TeamMember{}, fmt.Errorf("team member not found")
		}
		return entity.TeamMember{}, fmt.Errorf("failed to scan team member: %w", err)
	}

	return member, nil
}

func (r *teamRepository) AddTeamMember(ctx context.Context, member entity.TeamMember) error {
	query := sq.Insert("team_members").
		Columns("team_id", "user_id", "role", "joined_at").
		Values(member.TeamID, member.UserID, member.Role, member.JoinedAt).
		RunWith(r.db).
		PlaceholderFormat(sq.Question)

	_, err := query.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to add team member: %w", err)
	}

	return nil
}
