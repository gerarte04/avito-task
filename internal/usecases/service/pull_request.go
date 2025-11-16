package service

import (
	"avito-task/internal/domain"
	"avito-task/internal/repository"
	"avito-task/internal/usecases"
	"context"
	"errors"
	"fmt"

	"avito-task/pkg/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestService struct {
	pool *pgxpool.Pool
	prRepo repository.PullRequestRepo
	userRepo repository.UserRepo
}

func NewPullRequestService(
	pool *pgxpool.Pool,
	prRepo repository.PullRequestRepo,
	userRepo repository.UserRepo,
	) *PullRequestService {
	return &PullRequestService{
		pool: pool,
		prRepo: prRepo,
		userRepo: userRepo,
	}
}

func (s *PullRequestService) CreatePullRequest(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	const op = "PullRequestService.CreatePullRequest"

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin tx: %w", op, err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	author, err := s.userRepo.GetByID(ctx, tx, pr.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pr, err = s.prRepo.CreatePullRequest(ctx, tx, pr)
	if err != nil {
		if errors.Is(err, database.ErrUniqueViolation) {
			return nil, fmt.Errorf("%s: %w", op, usecases.ErrPRIDExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rews, err := s.userRepo.GetByTeam(ctx, tx, repository.GetByTeamOpts{
		TeamName: author.TeamName,
		OnlyActive: true,
		Limit: 2,
		ExcludeIDs: []string{author.ID},
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = s.prRepo.AddReviewers(ctx, tx, pr.ID, rews); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for _, r := range rews {
		pr.Reviewers = append(pr.Reviewers, r.ID)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit tx: %w", op, err)
	}

	return pr, nil
}

func (s *PullRequestService) Merge(ctx context.Context, id string) (*domain.PullRequest, error) {
	const op = "PullRequestService.Merge"

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: failed to begin tx: %w", op, err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	pr, err := s.prRepo.Merge(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	pr.Reviewers, err = s.prRepo.GetReviewers(ctx, tx, pr.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: failed to commit tx: %w", op, err)
	}

	return pr, nil
}

func (s *PullRequestService) Reassign(ctx context.Context, prID string, userID string) (string, *domain.PullRequest, error) {
	const op = "PullRequestService.Reassign"

	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})

	if err != nil {
		return "", nil, fmt.Errorf("%s: failed to begin tx: %w", op, err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	prev, err := s.userRepo.GetByID(ctx, tx, userID)
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	pr, err := s.prRepo.GetByID(ctx, tx, prID)
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	if pr.Status == domain.PRMerged {
		return "", nil, fmt.Errorf("%s: %w", op, usecases.ErrPRMerged)
	}

	rews, err := s.userRepo.GetByTeam(ctx, tx, repository.GetByTeamOpts{
		TeamName: prev.TeamName,
		OnlyActive: true,
		Limit: 1,
		ExcludeIDs: []string{prev.ID, pr.AuthorID},
	})
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(rews) == 0 {
		return "", nil, fmt.Errorf("%s: %w", op, usecases.ErrNoCandidate)
	}

	if err = s.prRepo.Reassign(ctx, tx, pr.ID, prev.ID, rews[0].ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil, fmt.Errorf("%s: %w", op, usecases.ErrNotAssigned)
		}

		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	pr.Reviewers, err = s.prRepo.GetReviewers(ctx, tx, pr.ID)
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return "", nil, fmt.Errorf("%s: failed to commit tx: %w", op, err)
	}

	return rews[0].ID, pr, nil
}
