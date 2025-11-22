package pr_service

import (
	"PRService/internal/model"
	"context"
)

func (s *Service) MergePullRequest(ctx context.Context, id model.PullRequestID) (*model.PullRequest, error) {
	_, err := s.storage.GetPullRequest(ctx, id)
	if err != nil {
		return nil, err
	}

	request, err := s.storage.MergePullRequest(ctx, id)
	if err != nil {
		return nil, err
	}

	return request, nil
}
