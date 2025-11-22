package router

import (
	"PRService/internal/errors"
	"PRService/internal/model"
	"github.com/gofiber/fiber/v3"
)

func (r *Router) CreatePullRequest(c fiber.Ctx) error {
	var body model.PullRequest
	if err := c.Bind().Body(&body); err != nil || !body.Valid() {
		r.logger.Warn(err)
		return c.
			Status(fiber.StatusBadRequest).
			JSON(err)
	}

	err := r.service.CreatePullRequest(c.Context(), &body)
	if err != nil {
		return c.Status(500).JSON(err)
	}
}
