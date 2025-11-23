package prservice

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
	"errors"
)

func (s *Service) CreatePullRequest(ctx context.Context, request *model.PullRequest) error {
	_, err := s.storage.GetPullRequest(ctx, request.ID)
	if err != nil && !errors.Is(err, service.ErrResourceNotFound) {
		s.logger.Errorf("CreatePullRequest: could not get pull request: %v", err)
		return err
	}
	if err == nil {
		return service.ErrPullRequestExists
	}

	_, err = s.storage.GetUser(ctx, request.AuthorID)
	if err != nil {
		s.logger.Errorf("CreatePullRequest: could not get user: %v", err)
		return err
	}

	request.Init()
	err = s.storage.CreatePullRequest(ctx, request)
	if err != nil {
		s.logger.Errorf("CreatePullRequest: could not create pull request: %v", err)
		return err
	}

	return nil
}
