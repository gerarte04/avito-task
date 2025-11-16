package repository

import (
	"avito-task/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type TeamRepo interface {
	CreateTeam(ctx context.Context, tx pgx.Tx, team *domain.Team) (bool, error)
	GetByName(ctx context.Context, tx pgx.Tx, name string) (*domain.Team, error)
}
