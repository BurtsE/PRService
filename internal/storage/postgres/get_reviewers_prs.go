package postgres

import (
	"PRService/internal/model"
	"context"
)

func (r *Repository) GetReviewersPRs(ctx context.Context, id model.UserID) ([]model.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}
