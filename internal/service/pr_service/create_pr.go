package pr_service

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
)

func (s *Service) CreatePullRequest(ctx context.Context, request *model.PullRequest) error {
	exists, err := s.storage.PullRequestExists(ctx, request.Id)
	if err != nil {
		return err
	}
	if exists {
		return service.ErrPullRequestExists
	}

	exists, err = s.storage.UserExists(ctx, request.AuthorId)
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
