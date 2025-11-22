package postgres

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
)

func (r *Repository) CreatePullRequest(ctx context.Context, request *model.PullRequest) error {
	query := `
		INSERT INTO pull_requests (id, name, author_id, status, created_at, merged_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	tx, err := r.c.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, query, request.Id, request.Name, request.AuthorId, request.Status, request.CreatedAt, request.MergedAt)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO pull_request_reviewers(pull_request_id, user_id)
			SELECT $1, id
			FROM users
			where
			    team_name = (SELECT team_name FROM users where id = $2)
				AND is_active = true
				AND id != $2
			ORDER BY RANDOM()
			LIMIT 2
		RETURNING user_id
	`

	rows, err := tx.Query(ctx, query, request.Id, request.AuthorId)
	if err != nil {
		return err
	}
	defer rows.Close()

	var userId model.UserID
	for rows.Next() {
		if err := rows.Scan(&userId); err != nil {
			return err
		}
		request.Reviewers = append(request.Reviewers, userId)
	}
	if len(request.Reviewers) == 0 {
		return service.ErrReviewerNotAssigned
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
