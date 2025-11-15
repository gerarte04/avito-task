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

func (r *PullRequestRepo) CreatePullRequest(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	const op = "PullRequestRepo.CreatePullRequest"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	sql := `
		INSERT INTO pull_requests (id, name, author_id)
		VALUES ($1, $2, $3)
		RETURNING status, created_at`

	if err = tx.QueryRow(
		ctx, sql, pr.ID, pr.Name, pr.AuthorID,
	).Scan(&pr.Status, &pr.CreatedAt); err != nil {
		dbErr := pkgPostgres.DetectError(err)

		if errors.Is(dbErr, database.ErrUniqueViolation) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrPRIDExists)
		} else if errors.Is(dbErr, database.ErrForeignKeyViolation) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrUserNotExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	sql = `
		WITH assigned_rews AS (
			SELECT id FROM users
			WHERE team_name = (SELECT team_name FROM users WHERE id = $1)
			AND is_active = TRUE AND id != $1
			LIMIT 2
		)
		INSERT INTO reviewers (pr_id, user_id)
		SELECT $2, id
		FROM assigned_rews
		RETURNING user_id`

	rows, err := tx.Query(ctx, sql, pr.AuthorID, pr.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var rev string

		if err = rows.Scan(&rev); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		pr.Reviewers = append(pr.Reviewers, rev)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return pr, nil
}

func (r *PullRequestRepo) Merge(ctx context.Context, id string) (*domain.PullRequest, error) {
	const op = "PullRequestRepo.Merge"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	sql := `
		UPDATE pull_requests SET
		status = 'MERGED',
		merged_at = COALESCE(merged_at, CURRENT_TIMESTAMP)
		WHERE id = $1
		RETURNING id, name, author_id, status, created_at, merged_at`

	var pr domain.PullRequest
	if err = tx.QueryRow(ctx, sql, id).Scan(
		&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrPRNotExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	sql = "SELECT user_id FROM reviewers WHERE pr_id = $1"

	rows, err := tx.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var userID string

		if err = rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		pr.Reviewers = append(pr.Reviewers, userID)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &pr, nil
}

func (r *PullRequestRepo) Reassign(ctx context.Context, prID string, userID string) (string, *domain.PullRequest, error) {
	const op = "PullRequestRepo.Reassign"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", op, err)
	}
	
	defer func() { _ = tx.Rollback(ctx) }()

	sql := "SELECT team_name FROM users WHERE id = $1"

	var teamName string
	if err = tx.QueryRow(ctx, sql, userID).Scan(&teamName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil, fmt.Errorf("%s: %w", op, repository.ErrUserNotExists)
		}

		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	sql = `
		SELECT id, name, author_id, status, created_at
		FROM pull_requests WHERE id = $1`

	var pr domain.PullRequest
	if err = tx.QueryRow(ctx, sql, prID).Scan(
		&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil, fmt.Errorf("%s: %w", op, repository.ErrPRNotExists)
		}

		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	if pr.Status == domain.PRMerged {
		return "", nil, fmt.Errorf("%s: %w", op, repository.ErrPRMerged)
	}

	sql = `
		WITH assigned_rews AS (
			SELECT id FROM users
			WHERE team_name = $1 AND is_active = TRUE
			AND id != $2 AND id != $3
			LIMIT 1
		)
		UPDATE reviewers SET
		user_id = (SELECT id FROM assigned_rews)
		WHERE pr_id = $4 AND user_id = $2
		RETURNING user_id`

	var newUserID string
	if err = tx.QueryRow(
		ctx, sql, teamName, userID, pr.AuthorID, pr.ID,
	).Scan(&newUserID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil, fmt.Errorf("%s: %w", op, repository.ErrNotAssigned)
		}
		
		if errors.Is(pkgPostgres.DetectError(err), database.ErrNotNullViolation) {
			return "", nil, fmt.Errorf("%s: %w", op, repository.ErrNoCandidate)
		}

		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	sql = "SELECT user_id FROM reviewers WHERE pr_id = $1"

	rows, err := tx.Query(ctx, sql, prID)
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var rewId string

		if err = rows.Scan(&rewId); err != nil {
			return "", nil, fmt.Errorf("%s: %w", op, err)
		}

		pr.Reviewers = append(pr.Reviewers, rewId)
	}

	if err = tx.Commit(ctx); err != nil {
		return "", nil, fmt.Errorf("%s: %w", op, err)
	}

	return newUserID, &pr, nil
}
