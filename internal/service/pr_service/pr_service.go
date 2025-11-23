package prservice

import (
	"PRService/internal/service"
	"PRService/internal/storage"

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
