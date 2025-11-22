package router

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"github.com/gofiber/fiber/v3"
)

func (r *Router) getTeam(c fiber.Ctx) error {
	name := c.Query("team_name")
	if name == "" {
		return r.ProcessError(c, service.ErrResourceNotFound)
	}
	team, err := r.service.GetTeam(c.Context(), model.TeamName(name))
	if err != nil {
		return c.
			Status(fiber.StatusNotFound).
			JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(team)
}
