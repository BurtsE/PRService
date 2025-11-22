package postgres

import (
	"PRService/internal/model"
	"context"
)

func (r *Repository) PullRequestExists(ctx context.Context, id model.PullRequestID) (bool, error) {
	var exists bool
	err := r.c.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
