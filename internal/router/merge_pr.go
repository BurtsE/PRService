package router

import (
	"PRService/internal/model"
	"PRService/internal/service"

	"github.com/gofiber/fiber/v3"
)

func (r *Router) mergePullRequest(c fiber.Ctx) error {
	var body struct {
		RequestID string `json:"pull_request_id"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return r.ProcessError(c, err)
	}

	if body.RequestID == "" {
		return r.ProcessError(c, service.ErrResourceNotFound)
	}

	request, err := r.service.MergePullRequest(c.Context(), model.PullRequestID(body.RequestID))
	if err != nil {
		return r.ProcessError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(request)
}
