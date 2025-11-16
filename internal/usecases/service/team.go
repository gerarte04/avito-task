package service

import (
	"avito-task/internal/domain"
	"avito-task/internal/repository"
	"avito-task/internal/usecases"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamService struct {
	pool *pgxpool.Pool
	teamRepo repository.TeamRepo
	userRepo repository.UserRepo
}

func NewTeamService(
	pool *pgxpool.Pool,
	teamRepo repository.TeamRepo,
	userRepo repository.UserRepo,
) *TeamService {
	return &TeamService{
		pool: pool,
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	const op = "TeamService.CreateTeam"

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin tx: %w", op, err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	exists, err := s.teamRepo.CreateTeam(ctx, tx, team)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if exists {
		return nil, fmt.Errorf("%s: %w", op, usecases.ErrTeamNameExists)
	}

	for _, u := range team.Members {
		u.TeamName = team.Name
	}

	if err = s.userRepo.UpsertUsers(ctx, tx, team.Members); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit tx: %w", op, err)
	}

	return team, nil
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	const op = "TeamService.GetTeam"

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin tx: %w", op, err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	team, err := s.teamRepo.GetByName(ctx, tx, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	team.Members, err = s.userRepo.GetByTeam(ctx, tx, repository.GetByTeamOpts{TeamName: team.Name})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit tx: %w", op, err)
	}

	return team, nil
}
