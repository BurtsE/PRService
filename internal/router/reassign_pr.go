package router

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"github.com/gofiber/fiber/v3"
)

func (r *Router) reassignPullRequestReviewer(c fiber.Ctx) error {
	var body struct {
		PullRequestID model.PullRequestID `json:"pull_request_id"`
		UserID        model.UserID        `json:"old_reviewer_id"`
	}
	if err := c.Bind().Body(&body); err != nil {
		return r.ProcessError(c, err)
	}

	if body.PullRequestID == "" || body.UserID == "" {
		return r.ProcessError(c, service.ErrResourceNotFound)
	}

	pullRequest, err := r.service.ReassignPullRequestReviewer(c.Context(), body.PullRequestID, body.UserID)
	if err != nil {
		return r.ProcessError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(pullRequest)
}
