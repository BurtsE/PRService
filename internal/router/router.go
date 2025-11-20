package router

import (
	"PRService/internal/config"
	"PRService/internal/service"
	"context"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type Router struct {
	app     *fiber.App
	service service.Service
	logger  *logrus.Logger
}

func NewRouter(cfg *config.Config, logger *logrus.Logger, service service.Service) *Router {
	app := fiber.New(fiber.Config{})
	app.Server().Logger = logger
	return &Router{
		app:     app,
		service: service,
		logger:  logger,
	}
}

func (r *Router) Start(_ context.Context) error {
	return r.app.Server().ListenAndServe(":8080")
}
