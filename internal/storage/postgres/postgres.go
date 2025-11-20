package postgres

import (
	"PRService/internal/config"
	"PRService/internal/service"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

var _ service.Service = (*Repository)(nil)

type Repository struct {
	c *pgx.Conn
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	DSN := fmt.Sprintf("host = %s dbname=%s user=%s password=%s  sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Database,
		cfg.Postgres.User,
		cfg.Postgres.Password,
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
