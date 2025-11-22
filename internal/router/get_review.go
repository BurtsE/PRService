package router

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"github.com/gofiber/fiber/v3"
)

func (r *Router) getReview(c fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return r.ProcessError(c, service.ErrResourceNotFound)
	}

	pullRequests, err := r.service.GetReviewersPRs(c.Context(), model.UserID(userID))
	if err != nil {
		return r.ProcessError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"user_id":      userID,
		"pullRequests": pullRequests,
	})
}
