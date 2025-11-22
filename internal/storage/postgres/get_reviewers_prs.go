package postgres

import (
	"PRService/internal/model"
	"context"
)

func (r *Repository) GetReviewersPRs(ctx context.Context, userID model.UserID) ([]model.PullRequest, error) {
	var prIDs []model.PullRequestID
	query := `
		SELECT request_id
		FROM pull_request_reviewers
		WHERE reviewer_id = $1
	`
	rows, err := r.c.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prID model.PullRequestID
	for rows.Next() {
		if err = rows.Scan(&prID); err != nil {
			return nil, err
		}
		prIDs = append(prIDs, prID)
	}

	pullRequests := make([]model.PullRequest, len(prIDs))
	for index, id := range prIDs {
		pullRequests[index], err = r.GetPullRequest(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	return pullRequests, nil
}
