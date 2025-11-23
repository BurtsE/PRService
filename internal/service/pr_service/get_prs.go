package prservice

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
)

func (s *Service) GetReviewersPRs(ctx context.Context, userID model.UserID) ([]model.PullRequest, error) {
	exists, err := s.storage.UserExists(ctx, userID)
	if err != nil {
		s.logger.Errorf("Error checking if user exists: %v", err)
		return nil, err
	}
	if !exists {
		return nil, service.ErrResourceNotFound
	}

	pullRequests, err := s.storage.GetReviewersPRs(ctx, userID)
	if err != nil {
		s.logger.Errorf("Error getting reviewers PRs: %v", err)
		return nil, err
	}

	return pullRequests, nil
}
