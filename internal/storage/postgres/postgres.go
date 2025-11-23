package postgres

import (
	"PRService/internal/config"
	"PRService/internal/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ storage.Storage = (*Repository)(nil)

type Repository struct {
	c *pgxpool.Pool
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	// Use a URL format for the connection string with pgxpool
	DSN := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Database,
	)
	c, err := pgxpool.New(context.Background(), DSN)
	if err != nil {
		return nil, err
	}

	err = c.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return &Repository{
		c: c,
	}, nil
}

func (r *Repository) Close(ctx context.Context) error {
	r.c.Close()
	return nil
}
