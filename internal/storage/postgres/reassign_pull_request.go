package postgres

import (
	"PRService/internal/model"
	"context"
)

func (r *Repository) ReassignPullRequestReviewer(ctx context.Context,
	pullRequest *model.PullRequest, oldReviewerID model.UserID, newReviewerID model.UserID) error {

	query := `
		UPDATE pull_request_reviewers
		SET user_id = $1
		WHERE pull_request_id = $1
			AND user_id = $2
	`

	_, err := r.c.Exec(ctx, query, pullRequest.ID, oldReviewerID, newReviewerID)
	if err != nil {
		return err
	}

	return nil
}
