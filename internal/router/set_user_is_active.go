package router

import (
	"PRService/internal/errors"
	"PRService/internal/model"
	"github.com/gofiber/fiber/v3"
)

func (r *Router) SetUserIsActive(c fiber.Ctx) error {
	var body model.User
	if err := c.Bind().Body(&body); err != nil || !body.Valid() {
		r.logger.Warn(err)
		return c.
			Status(fiber.StatusBadRequest).
			JSON(errors.NewErrorResponse(errors.ResourceNotFound))
	}

	err := r.service.SetUserIsActive(c.Context(), &body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	return c.Status(200).JSON(body)
}
