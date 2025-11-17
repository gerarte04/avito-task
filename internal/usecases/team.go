package usecases

import (
	"avito-task/internal/domain"
	"context"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error)
	GetTeam(ctx context.Context, name string) (*domain.Team, error)
	GetTeamStats(ctx context.Context, name string) (*domain.TeamStats, error)
	DeactivateTeam(ctx context.Context, name string) ([]*domain.User, error)
}
