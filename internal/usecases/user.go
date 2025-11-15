package usecases

import (
	"avito-task/internal/domain"
	"context"
)

type UserService interface {
	SetIsActive(ctx context.Context, id string, isActive bool) (*domain.User, error)
	GetReview(ctx context.Context, id string) ([]*domain.PullRequestShort, error)
}
