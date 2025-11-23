package prservice

import (
	"PRService/internal/model"
	"context"
)

func (s *Service) GetReviewersPRs(ctx context.Context, userID model.UserID) ([]model.PullRequest, error) {
	_, err := s.storage.GetUser(ctx, userID)
	if err != nil {
		s.logger.Errorf("GetReviewersPRs: could not get user: %v", err)
		return nil, err
	}

	pullRequests, err := s.storage.GetReviewersPRs(ctx, userID)
	if err != nil {
		s.logger.Errorf("GetReviewersPRs: could not get reviewers PRs: %v", err)
		return nil, err
	}

	return pullRequests, nil
}
