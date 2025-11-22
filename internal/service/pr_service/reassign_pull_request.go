package pr_service

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
)

func (s *Service) ReassignPullRequestReviewer(ctx context.Context,
	pullRequestID model.PullRequestID, userID model.UserID) (*model.PullRequest, error) {

	exists, err := s.storage.PullRequestExists(ctx, pullRequestID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, service.ErrResourceNotFound
	}

	pullRequest, err := s.storage.ReassignPullRequestReviewer(ctx, pullRequestID, userID)
	if err != nil {
		return nil, err
	}

	return pullRequest, nil
}
