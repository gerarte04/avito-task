package service

import (
	"avito-task/internal/domain"
	"avito-task/internal/repository"
	"context"
	"fmt"
)

type PullRequestService struct {
	prRepo repository.PullRequestRepo
}

func NewPullRequestService(repo repository.PullRequestRepo) *PullRequestService {
	return &PullRequestService{
		prRepo: repo,
	}
}

func (s *PullRequestService) CreatePullRequest(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	const op = "PullRequestService.CreatePullRequest"

	res, err := s.prRepo.CreatePullRequest(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (s *PullRequestService) Merge(ctx context.Context, id string) (*domain.PullRequest, error) {
	const op = "PullRequestService.Merge"

	pr, err := s.prRepo.Merge(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pr, nil
}

func (s *PullRequestService) Reassign(ctx context.Context, prID string, userID string) (string, *domain.PullRequest, error) {
	const op = "PullRequestService.Reassign"

	newUserID, pr, err := s.prRepo.Reassign(ctx, prID, userID)
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	return newUserID, pr, err
}
