package router

import (
	"PRService/internal/errors"
	"PRService/internal/model"
	"github.com/gofiber/fiber/v3"
	"net/http"
)

/*
AddTeam creates or updates team
Users are being created if not exist
*/
func (r *Router) AddTeam(c fiber.Ctx) error {
	var body *model.Team
	if err := c.Bind().Body(body); err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON(errors.NewErrorResponse(errors.ResourceNotFound))
	}

	err := r.service.CreateTeam(body)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	return c.Status(http.StatusCreated).JSON(body)
}
