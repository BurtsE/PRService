package pr_service

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"PRService/internal/storage"
	"context"
	"github.com/sirupsen/logrus"
)

var _ service.Service = (*Service)(nil)

type Service struct {
	logger  *logrus.Logger
	storage storage.Storage
}

func NewService(logger *logrus.Logger, storage storage.Storage) *Service {
	return &Service{
		logger:  logger,
		storage: storage,
	}
}

func (s *Service) GetReviewersPRs(ctx context.Context, id model.UserID) ([]model.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}
