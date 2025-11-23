package prservice

import (
	"PRService/internal/model"
	"context"
)

func (s *Service) MergePullRequest(ctx context.Context, id model.PullRequestID) (*model.PullRequest, error) {
	_, err := s.storage.GetPullRequest(ctx, id)
	if err != nil {
		s.logger.Errorf("MergePullRequest: could not get pull request: %v", err)
		return nil, err
	}

	request, err := s.storage.MergePullRequest(ctx, id)
	if err != nil {
		s.logger.Errorf("MergePullRequest: could not merge pull request: %v", err)
		return nil, err
	}

	return request, nil
}
