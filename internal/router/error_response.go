package router

import (
	"PRService/internal/service"
	"errors"
	"github.com/gofiber/fiber/v3"
)

type ErrorCode string

const (
	TeamExists           ErrorCode = "TEAM_EXISTS"
	PullRequestExists    ErrorCode = "PR_EXISTS"
	PullRequestMerged    ErrorCode = "PR_MERGED"
	ReviewerNotAssigned  ErrorCode = "NOT_ASSIGNED"
	ReviewersUnavailable ErrorCode = "NO_CANDIDATE"
	ResourceNotFound     ErrorCode = "NOT_FOUND"
	InternalServerError  ErrorCode = "INTERNAL_SERVER_ERROR"
)

const internalServerErrorMessage = "unexpected error occured"

type ErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (r *Router) ProcessError(c fiber.Ctx, err error) error {
	var response ErrorResponse
	response.Message = err.Error()

	if errors.Is(err, service.ErrTeamExists) {
		response.Code = TeamExists
		c.Status(fiber.StatusBadRequest)
	} else if errors.Is(err, service.ErrPullRequestExists) {
		response.Code = PullRequestExists
		c.Status(fiber.StatusConflict)
	} else if errors.Is(err, service.ErrPullRequestMerged) {
		response.Code = PullRequestMerged
		c.Status(fiber.StatusBadRequest)
	} else if errors.Is(err, service.ErrReviewerNotAssigned) {
		response.Code = ReviewerNotAssigned
		c.Status(fiber.StatusConflict)
	} else if errors.Is(err, service.ErrReviewersUnavailable) {
		response.Code = ReviewersUnavailable
		c.Status(fiber.StatusBadRequest)
	} else if errors.Is(err, service.ErrResourceNotFound) {
		response.Code = ResourceNotFound
		c.Status(fiber.StatusNotFound)
	} else {
		response.Code = InternalServerError
		response.Message = internalServerErrorMessage
		c.Status(fiber.StatusInternalServerError)
	}

	return c.JSON(response)
}
