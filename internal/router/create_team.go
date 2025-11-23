package router

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"github.com/gofiber/fiber/v3"
	"net/http"
)

/*
AddTeam creates or updates team
Users are being created if not exist
*/
func (r *Router) createTeam(c fiber.Ctx) error {
	var team model.Team
	if err := c.Bind().Body(&team); err != nil {
		return r.ProcessError(c, err)
	}

	if !team.Valid() {
		return r.ProcessError(c, service.ErrResourceNotFound)
	}

	r.logger.Debugf("AddTeam %+v", team)
	err := r.service.CreateTeam(c.Context(), &team)
	if err != nil {
		return r.ProcessError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(map[string]any{
		"team": team,
	})
}
