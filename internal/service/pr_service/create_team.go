package pr_service

import (
	"PRService/internal/model"
	"context"
)

func (s *Service) CreateTeam(ctx context.Context, team *model.Team) error {
	err := s.storage.CreateTeam(ctx, team)
	if err != nil {
		return err
	}
	return nil
}
