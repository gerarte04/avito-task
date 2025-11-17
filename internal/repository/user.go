package repository

import (
	"avito-task/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type GetByTeamOpts struct {
	TeamName   string
	OnlyActive bool
	Limit      int
	ExcludeIDs []string
}

type UserRepo interface {
	GetByID(ctx context.Context, tx pgx.Tx, id string) (*domain.User, error)
	GetByTeam(ctx context.Context, tx pgx.Tx, opts GetByTeamOpts) ([]*domain.User, error)
	SetIsActive(ctx context.Context, id string, isActive bool) (*domain.User, error)
	DeactivateTeam(ctx context.Context, tx pgx.Tx, teamName string) ([]*domain.User, error)
	UpsertUsers(ctx context.Context, tx pgx.Tx, users []*domain.User) error
}
