package router

import (
	"PRService/internal/errors"
	"PRService/internal/model"
	"github.com/gofiber/fiber/v3"
)

func (r *Router) getTeam(c fiber.Ctx) error {
	name := c.Query("team_name")
	if name == "" {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(errors.NewErrorResponse(errors.ResourceNotFound))
	}
	team, err := r.service.GetTeam(c.Context(), model.TeamName(name))
	if err != nil {
		return c.JSON(err)
	}

	return c.JSON(team)
}
