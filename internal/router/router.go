package router

import (
	"PRService/internal/config"
	"PRService/internal/service"
	"context"
	"github.com/gofiber/fiber/v3"
	fiberlogger "github.com/gofiber/fiber/v3/middleware/logger"
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
	app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "[${time}] ${ip} ${method} ${path} ${status} - ${latency}\n",
	}))
	return &Router{
		app:     app,
		service: service,
		logger:  logger,
	}
}

func (r *Router) SetupRoutes() {
	r.app.Post("/team/add", r.CreateTeam)
	r.app.Get("/team/get", r.getTeam)

	r.app.Post("/users/setIsActive", r.SetUserIsActive)

	r.app.Get("/ping", func(c fiber.Ctx) error {
		return c.JSON("pong")
	})
}

func (r *Router) Start(_ context.Context, addr string) error {
	r.SetupRoutes()
	return r.app.Listen(addr)
}

func (r *Router) Stop(_ context.Context) error {
	return r.app.Server().Shutdown()
}
