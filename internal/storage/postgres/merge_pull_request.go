package postgres

import (
	"PRService/internal/model"
	"context"
)

func (r *Repository) MergePullRequest(ctx context.Context, id model.PullRequestID) (*model.PullRequest, error) {
	query := `
		UPDATE pull_requests
		SET status = $2 merged_at = NOW()
		WHERE id = $1
		RETURNING id, name, author_id, status, created_at, merged_at
	`
	var pr model.PullRequest
	err := r.c.QueryRow(ctx, query, id, model.PullRequestStatusMerged).Scan(
		&pr.Id,
		&pr.Name,
		&pr.AuthorId,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		return nil, err
	}

	return &pr, nil

}
