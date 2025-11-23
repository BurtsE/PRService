package postgres

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
	"database/sql"
	"errors"
)

func (r *Repository) GetPullRequest(ctx context.Context, pullRequestID model.PullRequestID) (model.PullRequest, error) {
	query := `
		SELECT id, name, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE id = $1
	`
	var pullRequest model.PullRequest
	err := r.c.QueryRow(ctx, query, pullRequestID).Scan(
		&pullRequest.ID,
		&pullRequest.Name,
		&pullRequest.AuthorID,
		&pullRequest.Status,
		&pullRequest.CreatedAt,
		&pullRequest.MergedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.PullRequest{}, service.ErrResourceNotFound
	}
	if err != nil {
		return model.PullRequest{}, err
	}

	query = `
		SELECT user_id
		FROM pull_request_reviewers
		WHERE pull_request_id = $1
	`

	rows, err := r.c.Query(ctx, query, pullRequestID)
	if err != nil {
		return model.PullRequest{}, err
	}
	defer rows.Close()

	var userID model.UserID
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return model.PullRequest{}, err
		}
		pullRequest.Reviewers = append(pullRequest.Reviewers, userID)
	}

	return pullRequest, nil
}
