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

	result := make([]PullRequestDto, 0, len(pullRequests))
	for i := range pullRequests {
		result = append(result, PullRequestDto{
			PullRequestID:   pullRequests[i].ID,
			PullRequestName: pullRequests[i].Name,
			AuthorID:        pullRequests[i].AuthorID,
			Status:          pullRequests[i].Status,
		})
	}
	response := map[string]any{
		"user_id":      userID,
		"pullRequests": result,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

type PullRequestDto struct {
	PullRequestID   model.PullRequestID     `json:"pull_request_id"`
	PullRequestName string                  `json:"pull_request_name"`
	AuthorID        model.UserID            `json:"author_id"`
	Status          model.PullRequestStatus `json:"status"`
}
