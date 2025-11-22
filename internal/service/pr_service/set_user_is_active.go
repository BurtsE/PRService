package pr_service

import (
	"PRService/internal/errors"
	"PRService/internal/model"
	"context"
)

func (s *Service) SetUserIsActive(ctx context.Context, user *model.User) error {
	exists, err := s.storage.UserExists(ctx, user.Id)
	if err != nil {
		s.logger.Errorf("Error checking if user exists: %v", err)
		return errors.NewErrorResponse(errors.InternalServerError)
	}
	if !exists {
		return errors.NewErrorResponse(errors.ResourceNotFound)
	}
	err = s.storage.SetUserIsActive(ctx, user)
	if err != nil {
		s.logger.Errorf("Error setting user isActive: %v", err)
		return errors.NewErrorResponse(errors.InternalServerError)
	}

	return nil
}
