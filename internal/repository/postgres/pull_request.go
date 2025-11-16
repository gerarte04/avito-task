package postgres

import (
	"avito-task/internal/domain"
	"avito-task/internal/repository"
	"context"
	"errors"
	"fmt"

	"avito-task/pkg/database"
	pkgPostgres "avito-task/pkg/database/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepo struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepo(pool *pgxpool.Pool) *PullRequestRepo {
	return &PullRequestRepo{
		pool: pool,
	}
}

// Reviews ---------------------------------------------------------

func (r *PullRequestRepo) GetReviewers(ctx context.Context, tx pgx.Tx, prID string) ([]string, error) {
	const op = "PullRequestRepo.GetReviewers"

	sql := "SELECT user_id FROM reviewers WHERE pr_id = $1"

	rows, err := tx.Query(ctx, sql, prID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()
	rw := []string{}

	for rows.Next() {
		var userID string

		if err = rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		rw = append(rw, userID)
	}

	return rw, nil
}

func (r *PullRequestRepo) GetUserReviews(ctx context.Context, id string) ([]*domain.PullRequestShort, error) {
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

	defer rows.Close()
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

func (r *PullRequestRepo) AddReviewers(ctx context.Context, tx pgx.Tx, prID string, users []*domain.User) error {
	const op = "PullRequestRepo.AddReviewers"

	sql := "INSERT INTO reviewers (pr_id, user_id) VALUES %s"

	values := ""
	args := []any{}

	for i, u := range users {
		comma := ","
		if i == len(users) - 1 {
			comma = ""
		}

		idx := i * 2 + 1
		values += fmt.Sprintf("($%d, $%d)%s ", idx, idx + 1, comma)
		args = append(args, prID, u.ID)
	}

	sql = fmt.Sprintf(sql, values)

	if _, err := tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// PRs ------------------------------------------------------------

func (r *PullRequestRepo) GetByID(ctx context.Context, tx pgx.Tx, id string) (*domain.PullRequest, error) {
	const op = "PullRequestRepo.GetByID"

	sql := `
		SELECT id, name, author_id, status, created_at, COALESCE(merged_at, '0001-01-01'::date)
		FROM pull_requests WHERE id = $1`

	var pr domain.PullRequest
	if err := tx.QueryRow(ctx, sql, id).Scan(
		&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrPRNotExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &pr, nil
}

func (r *PullRequestRepo) CreatePullRequest(ctx context.Context, tx pgx.Tx, pr *domain.PullRequest) (*domain.PullRequest, error) {
	const op = "PullRequestRepo.CreatePullRequest"

	sql := `
		INSERT INTO pull_requests (id, name, author_id)
		VALUES ($1, $2, $3)
		RETURNING status, created_at`

	if err := tx.QueryRow(
		ctx, sql, pr.ID, pr.Name, pr.AuthorID,
	).Scan(&pr.Status, &pr.CreatedAt); err != nil {
		dbErr := pkgPostgres.DetectError(err)

		if errors.Is(dbErr, database.ErrUniqueViolation) {
			return nil, fmt.Errorf("%s: %w", op, database.ErrUniqueViolation)
		} else if errors.Is(dbErr, database.ErrForeignKeyViolation) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrUserNotExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pr, nil
}

func (r *PullRequestRepo) Merge(ctx context.Context, tx pgx.Tx, id string) (*domain.PullRequest, error) {
	const op = "PullRequestRepo.Merge"

	sql := `
		UPDATE pull_requests SET
		status = 'MERGED',
		merged_at = COALESCE(merged_at, CURRENT_TIMESTAMP)
		WHERE id = $1
		RETURNING id, name, author_id, status, created_at, merged_at`

	var pr domain.PullRequest
	if err := tx.QueryRow(ctx, sql, id).Scan(
		&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrPRNotExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if pr.MergedAt.IsZero() {
		pr.MergedAt = nil
	}

	return &pr, nil
}

func (r *PullRequestRepo) Reassign(ctx context.Context, tx pgx.Tx, prID string, prevID string, newID string) error {
	const op = "PullRequestRepo.Reassign"
	
	sql := `
		UPDATE reviewers SET user_id = $1
		WHERE pr_id = $2 AND user_id = $3
		RETURNING user_id`

	if err := tx.QueryRow(
		ctx, sql, newID, prID, prevID,
	).Scan(&newID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, pgx.ErrNoRows)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
