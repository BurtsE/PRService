package pr_service

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
)

func (s *Service) CreateTeam(ctx context.Context, team *model.Team) error {
	exist, err := s.storage.TeamExists(ctx, team.Name)
	if err != nil {
		s.logger.Warn(err)
		return err
	}
	if exist {
		return service.ErrTeamExists
	}
	err = s.storage.CreateTeam(ctx, team)
	if err != nil {
		s.logger.Warn(err)
		return err
	}

	return nil
}
