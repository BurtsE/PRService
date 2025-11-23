package router

import (
	"PRService/internal/model"
	"PRService/internal/service"

	"github.com/gofiber/fiber/v3"
)

func (r *Router) createPullRequest(c fiber.Ctx) error {
	var request model.PullRequest
	if err := c.Bind().Body(&request); err != nil {
		r.logger.Warn(err)
		return r.ProcessError(c, err)
	}

	if !request.Valid() {
		return r.ProcessError(c, service.ErrResourceNotFound)
	}

	err := r.service.CreatePullRequest(c.Context(), &request)
	if err != nil {
		return r.ProcessError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(request)
}
