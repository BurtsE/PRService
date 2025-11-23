package router

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"github.com/gofiber/fiber/v3"
)

func (r *Router) setUserIsActive(c fiber.Ctx) error {
	var body model.User
	if err := c.Bind().Body(&body); err != nil {
		r.logger.Warn(err)
		return r.ProcessError(c, err)
	}

	if !body.Valid() {
		return r.ProcessError(c, service.ErrResourceNotFound)
	}

	err := r.service.SetUserIsActive(c.Context(), &body)
	if err != nil {
		return r.ProcessError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(body)
}
