package pr_service

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
)

func (s *Service) MergePullRequest(ctx context.Context, id model.PullRequestID) (*model.PullRequest, error) {
	exists, err := s.storage.PullRequestExists(ctx, id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, service.ErrResourceNotFound
	}

	request, err := s.storage.MergePullRequest(ctx, id)
	if err != nil {
		return nil, err
	}

	return request, nil
}
