package service

import (
	"avito-task/internal/domain"
	"avito-task/internal/repository"
	"context"
	"fmt"
)

type TeamService struct {
	teamRepo repository.TeamRepo
}

func NewTeamService(repo repository.TeamRepo) *TeamService {
	return &TeamService{
		teamRepo: repo,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	const op = "TeamService.CreateTeam"

	res, err := s.teamRepo.CreateTeam(ctx, team)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	const op = "TeamService.GetTeam"

	team, err := s.teamRepo.GetTeam(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return team, nil
}
