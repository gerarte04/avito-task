package service

import (
	"avito-task/internal/domain"
	"avito-task/internal/repository"
	"context"
	"fmt"
)

type UserService struct {
	userRepo repository.UserRepo
	prRepo repository.PullRequestRepo
}

func NewUserService(
	userRepo repository.UserRepo,
	prRepo repository.PullRequestRepo,
) *UserService {
	return &UserService{
		userRepo: userRepo,
		prRepo: prRepo,
	}
}

func (s *UserService) SetIsActive(ctx context.Context, id string, isActive bool) (*domain.User, error) {
	const op = "UserService.SetIsActive"

	user, err := s.userRepo.SetIsActive(ctx, id, isActive)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UserService) GetReview(ctx context.Context, id string) ([]*domain.PullRequestShort, error) {
	const op = "UserService.GetReview"

	prs, err := s.prRepo.GetUserReviews(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return prs, nil
}
