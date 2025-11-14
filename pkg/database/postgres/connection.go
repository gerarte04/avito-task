package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host 				string			`yaml:"host" env:"POSTGRES_HOST" env-required:"true"`
	Port 				string			`yaml:"port" env:"POSTGRES_PORT" env-required:"true"`
	DbName 				string			`yaml:"db" env:"POSTGRES_DB" env-required:"true"`
	User 				string			`yaml:"user" env:"POSTGRES_USER" env-required:"true"`
	Password 			string			`yaml:"password" env:"POSTGRES_PASSWORD" env-required:"true"`
	ConnectionTimeout	time.Duration	`yaml:"connection_timeout" env-default:"300ms"`
}

func NewPostgresPool(cfg Config) (*pgxpool.Pool, error) {
	const method = "postgres.NewPostgresPool"

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectionTimeout)
	defer cancel()

	pginfo := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.DbName, cfg.User, cfg.Password,
	)

	pool, err := pgxpool.New(ctx, pginfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", method, err)
	}

	return pool, nil
}
