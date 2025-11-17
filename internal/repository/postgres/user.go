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

func (r *UserRepo) GetByID(ctx context.Context, tx pgx.Tx, id string) (*domain.User, error) {
	const op = "UserRepo.GetByID"

	sql := "SELECT id, name, team_name, is_active FROM users WHERE id = $1"

	var user domain.User
	if err := tx.QueryRow(ctx, sql, id).Scan(
		&user.ID, &user.Name, &user.TeamName, &user.IsActive,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrUserNotExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (r *UserRepo) GetByTeam(ctx context.Context, tx pgx.Tx, opts repository.GetByTeamOpts) ([]*domain.User, error) {
	const op = "UserRepo.GetByTeam"
	
	sql := "SELECT id, name, team_name, is_active FROM users WHERE team_name = $1"
	args := []any{opts.TeamName}
	i := 2

	if opts.OnlyActive {
		sql = fmt.Sprintf("%s AND is_active = TRUE", sql)
	}

	for _, e := range opts.ExcludeIDs {
		sql = fmt.Sprintf("%s AND id != $%d", sql, i)
		args = append(args, e)
		i++
	}

	if opts.Limit > 0 {
		sql = fmt.Sprintf("%s ORDER BY RANDOM() LIMIT $%d", sql, i)
		args = append(args, opts.Limit)
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()
	users := []*domain.User{}

	for rows.Next() {
		var u domain.User
		
		if err = rows.Scan(&u.ID, &u.Name, &u.TeamName, &u.IsActive); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		
		users = append(users, &u)
	}

	return users, nil
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

func (r *UserRepo) DeactivateTeam(ctx context.Context, tx pgx.Tx, teamName string) ([]*domain.User, error) {
	const op = "UserRepo.DeactivateTeam"

	sql := `
		UPDATE users SET is_active = FALSE WHERE team_name = $1
		RETURNING id, name, team_name, is_active`

	rows, err := tx.Query(ctx, sql, teamName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()
	users := []*domain.User{}

	for rows.Next() {
		var u domain.User

		if err = rows.Scan(&u.ID, &u.Name, &u.TeamName, &u.IsActive); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		users = append(users, &u)
	}

	return users, nil
}

func (r *UserRepo) UpsertUsers(ctx context.Context, tx pgx.Tx, users []*domain.User) error {
	const op = "UserRepo.UpsertUsers"
	
	sql := `
		INSERT INTO users (id, name, team_name, is_active)
		VALUES %s
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name, team_name = EXCLUDED.team_name, is_active = EXCLUDED.is_active`

	values := ""
	args := []any{}

	for i, u := range users {
		comma := ","
		if i == len(users) - 1 {
			comma = ""
		}

		idx := i * 4 + 1
		values += fmt.Sprintf("($%d, $%d, $%d, $%d)%s ", idx, idx + 1, idx + 2, idx + 3, comma)
		args = append(args, u.ID, u.Name, u.TeamName, u.IsActive)
	}

	sql = fmt.Sprintf(sql, values)

	if _, err := tx.Exec(ctx, sql, args...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
