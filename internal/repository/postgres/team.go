package postgres

import (
	"avito-task/internal/domain"
	"avito-task/internal/repository"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepo struct {
	pool *pgxpool.Pool
}

func NewTeamRepo(pool *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{
		pool: pool,
	}
}

func (r *TeamRepo) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	const op = "TeamRepo.CreateTeam"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	sql := `
		INSERT INTO teams (name) VALUES ($1)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING (xmax <> 0)`

	var wasExisting bool
	if err = tx.QueryRow(ctx, sql, team.Name).Scan(&wasExisting); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if wasExisting {
		return nil, fmt.Errorf("%s: %w", op, repository.ErrTeamNameExists)
	}

	sql = `
		INSERT INTO users (id, name, team_name, is_active)
		VALUES %s
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name, team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active`

	values := ""
	args := []any{}

	for i, u := range team.Members {
		comma := ","
		if i == len(team.Members) - 1 {
			comma = ""
		}

		idx := i * 4 + 1
		values += fmt.Sprintf("($%d, $%d, $%d, $%d)%s ", idx, idx + 1, idx + 2, idx + 3, comma)
		args = append(args, u.ID, u.Name, team.Name, u.IsActive)
	}

	sql = fmt.Sprintf(sql, values)

	if _, err = tx.Exec(ctx, sql, args...); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return team, nil
}

func (r *TeamRepo) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	const op = "TeamRepo.GetTeam"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	sql := "SELECT EXISTS(SELECT * FROM teams WHERE name = $1)"

	var exists bool
	if err = tx.QueryRow(ctx, sql, name).Scan(&exists); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !exists {
		return nil, fmt.Errorf("%s: %w", op, repository.ErrTeamNotExists)
	}

	sql = "SELECT id, name, is_active FROM users WHERE team_name = $1"
	
	rows, err := tx.Query(ctx, sql, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	team := domain.Team{Name: name}

	for rows.Next() {
		var member domain.User

		if err = rows.Scan(&member.ID, &member.Name, &member.IsActive); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		team.Members = append(team.Members, &member)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &team, nil
}
