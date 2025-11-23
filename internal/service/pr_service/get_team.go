package prservice

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
)

func (s *Service) GetTeam(ctx context.Context, teamName model.TeamName) (*model.Team, error) {
	exists, err := s.storage.TeamExists(ctx, teamName)
	if err != nil {
		s.logger.Errorf("GetTeam: could not check if team exists: %v", err)
		return nil, err
	}
	if !exists {
		return nil, service.ErrResourceNotFound
	}
	team, err := s.storage.GetTeam(ctx, teamName)
	if err != nil {
		s.logger.Errorf("GetTeam: could not get team: %v")
		return nil, err
	}

	return team, nil
}
