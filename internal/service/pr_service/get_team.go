package pr_service

import (
	"PRService/internal/errors"
	"PRService/internal/model"
	"context"
)

func (s *Service) GetTeam(ctx context.Context, teamName model.TeamName) (*model.Team, error) {
	exists, err := s.storage.TeamExists(ctx, teamName)
	if err != nil {
		s.logger.Errorf("Error checking if team exists: %v", err)
		return nil, errors.NewErrorResponse(errors.InternalServerError)
	}
	if !exists {
		return nil, errors.NewErrorResponse(errors.ResourceNotFound)
	}
	team, err := s.storage.GetTeam(ctx, teamName)
	if err != nil {
		s.logger.Errorf("Error getting team: %v", err)
		return nil, errors.NewErrorResponse(errors.InternalServerError)
	}

	return team, nil
}
