package postgres

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
)

func (r *Repository) ReassignPullRequestReviewer(ctx context.Context,
	pullRequest *model.PullRequest, oldReviewerID model.UserID) error {

	query := `
		UPDATE pull_request_reviewers as old
		SET user_id = COALESCE((
		    SELECT id 
		    FROM users
		    WHERE
		        team_name = (SELECT team_name FROM users where id = $3)
		    	AND is_active = true
				AND id != $3
				AND id NOT IN (
		    		SELECT user_id
		    		FROM pull_request_reviewers
		    		where pull_request_id = $1
				)
			ORDER BY RANDOM()
			LIMIT 1
		), old.user_id)
		WHERE old.pull_request_id = $1
			AND old.user_id = $2
		RETURNING user_id
	`

	var newReviewerID model.UserID
	err := r.c.QueryRow(ctx, query, pullRequest.ID, oldReviewerID, pullRequest.AuthorID).Scan(&newReviewerID)
	if err != nil {
		return err
	}
	if newReviewerID == oldReviewerID {
		return service.ErrReviewersUnavailable
	}

	for i := range pullRequest.Reviewers {
		if pullRequest.Reviewers[i] == oldReviewerID {
			pullRequest.Reviewers[i] = newReviewerID
		}
	}

	return nil
}
