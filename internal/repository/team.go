package repository

import (
	"avito-task/internal/domain"
	"context"
)

type TeamRepo interface {
	CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error)
	GetTeam(ctx context.Context, name string) (*domain.Team, error)
}
