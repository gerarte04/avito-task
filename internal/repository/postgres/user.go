package postgres

import (
	"avito-task/internal/domain"
	"avito-task/internal/repository"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		pool: pool,
	}
}

func (r *UserRepo) SetIsActive(ctx context.Context, id string, isActive bool) (*domain.User, error) {
	const op = "UserRepo.SetIsActive"

	sql := `
		UPDATE users SET is_active = $1 WHERE id = $2
		RETURNING id, name, team_name, is_active`

	row := r.pool.QueryRow(ctx, sql, isActive, id)
	var user domain.User

	if err := row.Scan(&user.ID, &user.Name, &user.TeamName, &user.IsActive); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrUserNotExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (r *UserRepo) GetReview(ctx context.Context, id string) ([]*domain.PullRequestShort, error) {
	const op = "UserRepo.GetReview"

	sql := `
		SELECT p.id, p.name, p.author_id, p.status
		FROM reviewers r
		JOIN pull_requests p ON r.pr_id = p.id
		WHERE r.user_id = $1`

	rows, err := r.pool.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var prs []*domain.PullRequestShort

	for rows.Next() {
		var pr domain.PullRequestShort

		if err = rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		prs = append(prs, &pr)
	}

	return prs, nil
}
