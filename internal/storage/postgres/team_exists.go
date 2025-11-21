package postgres

import (
	"PRService/internal/model"
	"context"
)

func (r *Repository) TeamExists(ctx context.Context, teamName model.TeamName) (bool, error) {
	var exists bool
	err := r.c.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1)`, teamName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
