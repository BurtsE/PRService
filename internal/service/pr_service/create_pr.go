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
		return err
	}
	if err == nil {
		return service.ErrPullRequestExists
	}

	_, err = s.storage.GetUser(ctx, request.AuthorID)
	if err != nil {
		return err
	}

	request.Init()
	err = s.storage.CreatePullRequest(ctx, request)
	if err != nil {
		return err
	}

	return nil
}
