package pr_service

import (
	"PRService/internal/errors"
	"PRService/internal/model"
	"context"
)

func (s *Service) CreateTeam(ctx context.Context, team *model.Team) error {
	exist, err := s.storage.TeamExists(ctx, team.Name)
	if err != nil {
		s.logger.Warn(err)
		return errors.NewErrorResponse(errors.InternalServerError)
	}
	if exist {
		return errors.NewErrorResponse(errors.TeamExists)
	}
	err = s.storage.CreateTeam(ctx, team)
	if err != nil {
		s.logger.Warn(err)
		return errors.NewErrorResponse(errors.InternalServerError)
	}
	return nil
}
