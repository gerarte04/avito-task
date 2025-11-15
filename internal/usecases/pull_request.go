package usecases

import (
	"avito-task/internal/domain"
	"context"
)

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
	Merge(ctx context.Context, id string) (*domain.PullRequest, error)
	Reassign(ctx context.Context, prID string, userID string) (string, *domain.PullRequest, error)
}
