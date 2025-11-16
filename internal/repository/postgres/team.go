package postgres

import (
	"avito-task/internal/domain"
	"avito-task/internal/repository"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (r *TeamRepo) CreateTeam(ctx context.Context, tx pgx.Tx, team *domain.Team) (bool, error) {
	const op = "TeamRepo.TryCreateTeam"

	sql := `
		INSERT INTO teams (name) VALUES ($1)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING (xmax <> 0)`

	var wasExisting bool
	if err := tx.QueryRow(ctx, sql, team.Name).Scan(&wasExisting); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return wasExisting, nil
}

func (r *TeamRepo) GetByName(ctx context.Context, tx pgx.Tx, name string) (*domain.Team, error) {
	const op = "TeamRepo.GetByName"

	sql := "SELECT name FROM teams WHERE name = $1"

	var team domain.Team
	if err := tx.QueryRow(ctx, sql, name).Scan(&team.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrTeamNotExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &team, nil
}
