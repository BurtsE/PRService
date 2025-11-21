package postgres

import (
	"PRService/internal/model"
	"context"
)

func (r *Repository) UserExists(ctx context.Context, id model.UserID) (bool, error) {
	var exists bool
	err := r.c.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
