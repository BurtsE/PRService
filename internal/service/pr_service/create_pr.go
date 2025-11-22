package pr_service

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

	exists, err := s.storage.UserExists(ctx, request.AuthorID)
	if err != nil {
		return err
	}
	if !exists {
		return service.ErrResourceNotFound
	}

	request.Init()
	err = s.storage.CreatePullRequest(ctx, request)
	if err != nil {
		return err
	}

	return nil
}
