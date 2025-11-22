package pr_service

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
)

func (s *Service) SetUserIsActive(ctx context.Context, user *model.User) error {
	exists, err := s.storage.UserExists(ctx, user.Id)
	if err != nil {
		s.logger.Errorf("Error checking if user exists: %v", err)
		return err
	}
	if !exists {
		return service.ErrResourceNotFound
	}

	err = s.storage.SetUserIsActive(ctx, user)
	if err != nil {
		s.logger.Errorf("Error setting user isActive: %v", err)
		return err
	}

	return nil
}
