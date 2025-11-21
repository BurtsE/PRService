package app

import (
	"PRService/internal/config"
	"PRService/internal/router"
	"PRService/internal/service/pr_service"
	"PRService/internal/storage/postgres"
	"context"
	"github.com/sirupsen/logrus"
)

type App struct {
	postgres *postgres.Repository
	service  *pr_service.Service
	logger   *logrus.Logger
	config   *config.Config
	router   *router.Router
}

func (a *App) Logger() *logrus.Logger {
	if a.logger == nil {
		a.logger = logrus.New()
	}
	return a.logger
}

func (a *App) Config() *config.Config {
	if a.config == nil {
		cfg, err := config.InitConfig()
		if err != nil {
			a.Logger().Panic(err)
		}
		a.config = cfg
	}
	return a.config
}

func (a *App) Postgres() *postgres.Repository {
	if a.postgres == nil {
		pg, err := postgres.NewRepository(a.Config())
		if err != nil {
			a.Logger().Panic(err)
		}
		a.postgres = pg
	}
	return a.postgres
}

func (a *App) Service() *pr_service.Service {
	if a.service == nil {
		a.service = pr_service.NewService(a.Logger(), a.Postgres())
	}
	return a.service
}

func (a *App) Router() *router.Router {
	if a.router == nil {
		a.router = router.NewRouter(a.Config(), a.Logger(), a.Service())
	}
	return a.router
}

func (a *App) Start(ctx context.Context) error {
	return a.Router().Start(ctx)
}

func (a *App) Shutdown(ctx context.Context) error {
	a.Logger().Info("closing databasse...")
	_ = a.postgres.Close(ctx)
	a.Logger().Info("shutting down server...")
	_ = a.Router().Stop(ctx)
	a.Logger().Info("app closed")
}
