package postgres

import (
	"PRService/internal/config"
	"PRService/internal/storage"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var _ storage.Storage = (*Repository)(nil)

type Repository struct {
	c *pgx.Conn
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	DSN := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		cfg.Host,
		cfg.Database,
		cfg.User,
		cfg.Password,
	)
	c, err := pgx.Connect(context.Background(), DSN)
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
	return r.c.Close(ctx)
}
