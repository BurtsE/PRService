package postgres

import (
	"PRService/internal/model"
	"context"
)

func (r *Repository) MergePullRequest(ctx context.Context, requestID model.PullRequestID) (*model.PullRequest, error) {
	query := `
		UPDATE pull_requests
		SET status = $2, merged_at = NOW()
		WHERE id = $1
		RETURNING id, name, author_id, status, created_at, merged_at
	`

	tx, err := r.c.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var pr model.PullRequest
	err = tx.QueryRow(ctx, query, requestID, model.PullRequestStatusMerged).Scan(
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

	query = `
		SELECT user_id
		FROM pull_request_reviewers
		WHERE pull_request_id = $1
	`
	rows, err := tx.Query(ctx, query, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID model.UserID
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		pr.Reviewers = append(pr.Reviewers, userID)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &pr, nil

}
