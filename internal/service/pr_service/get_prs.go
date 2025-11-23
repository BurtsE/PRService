package prservice

import (
	"PRService/internal/model"
	"context"
)

func (s *Service) GetReviewersPRs(ctx context.Context, userID model.UserID) ([]model.PullRequest, error) {
	_, err := s.storage.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	pullRequests, err := s.storage.GetReviewersPRs(ctx, userID)
	if err != nil {
		s.logger.Errorf("Error getting reviewers PRs: %v", err)
		return nil, err
	}

	return pullRequests, nil
}
