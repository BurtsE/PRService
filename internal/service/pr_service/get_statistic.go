package prservice

import (
	"PRService/internal/model"
	"context"
)

func (s *Service) GetStatistic(ctx context.Context) (model.Statistic, error) {
	return s.storage.GetStatistic(ctx)
}
