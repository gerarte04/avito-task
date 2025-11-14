package repository

import (
	"avito-task/internal/domain"
	"context"
)

type PullRequestRepo interface {
	CreatePullRequest(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error)
	Merge(ctx context.Context, id string) (*domain.PullRequest, error)
	Reassign(ctx context.Context, prID string, userID string) (string, *domain.PullRequest, error)
}
