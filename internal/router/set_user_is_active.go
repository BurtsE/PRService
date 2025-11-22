package router

import (
	"PRService/internal/model"
	"github.com/gofiber/fiber/v3"
)

func (r *Router) setUserIsActive(c fiber.Ctx) error {
	var body model.User
	if err := c.Bind().Body(&body); err != nil || !body.Valid() {
		r.logger.Warn(err)
		return r.ProcessError(c, err)
	}

	err := r.service.SetUserIsActive(c.Context(), &body)
	if err != nil {
		return r.ProcessError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(body)
}
