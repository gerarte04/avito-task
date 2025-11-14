package repository

import (
	"avito-task/internal/domain"
	"context"
)

type UserRepo interface {
	SetIsActive(ctx context.Context, id string, isActive bool) (*domain.User, error)
	GetReview(ctx context.Context, id string) ([]*domain.PullRequestShort, error)
}
