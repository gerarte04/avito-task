package repository

import (
	"avito-task/internal/domain"
	"context"

	"github.com/jackc/pgx/v5"
)

type PullRequestRepo interface {
	GetReviewers(ctx context.Context, tx pgx.Tx, prID string) ([]string, error)
	GetUserReviews(ctx context.Context, id string) ([]*domain.PullRequestShort, error)
	AddReviewers(ctx context.Context, tx pgx.Tx, prID string, users []*domain.User) error

	GetByID(ctx context.Context, tx pgx.Tx, id string) (*domain.PullRequest, error)
	CreatePullRequest(ctx context.Context, tx pgx.Tx, pr *domain.PullRequest) (*domain.PullRequest, error)
	Merge(ctx context.Context, tx pgx.Tx, id string) (*domain.PullRequest, error)
	Reassign(ctx context.Context, tx pgx.Tx, prID string, prevID string, newID string) error

	GetUserReviewsCounts(ctx context.Context, tx pgx.Tx, teamName string) ([]*domain.UserStats, error)
	GetPRReviewersCounts(ctx context.Context, tx pgx.Tx, teamName string) ([]*domain.PullRequestStats, error)
}
