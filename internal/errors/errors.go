package errors

import "fmt"

type ErrorCode string

const (
	TeamExists           ErrorCode = "TEAM_EXISTS"
	PullRequestExists    ErrorCode = "PR_EXISTS"
	PullRequestMerged    ErrorCode = "PR_MERGED"
	ReviewerNotAssigned  ErrorCode = "NOT_ASSIGNED"
	ReviewersUnavailable ErrorCode = "NO_CANDIDATE"
	ResourceNotFound     ErrorCode = "NOT_FOUND"
)

type ErrorResponse struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func NewErrorResponse(code ErrorCode) *ErrorResponse {
	switch code {
	case TeamExists:
		return &ErrorResponse{
			Code:    TeamExists,
			Message: "team already exists",
		}
	case PullRequestExists:
		return &ErrorResponse{
			Code:    PullRequestExists,
			Message: "pull request already exists",
		}
	case PullRequestMerged:
		return &ErrorResponse{
			Code:    PullRequestMerged,
			Message: "pull request already merged",
		}
	case ReviewerNotAssigned:
		return &ErrorResponse{
			Code:    ReviewerNotAssigned,
			Message: "reviewer not assigned",
		}
	case ReviewersUnavailable:
		return &ErrorResponse{
			Code:    ReviewersUnavailable,
			Message: "reviewers unavailable",
		}
	case ResourceNotFound:
		return &ErrorResponse{
			Code:    ResourceNotFound,
			Message: "resource not found",
		}
	default:
		return &ErrorResponse{
			Code:    code,
			Message: fmt.Sprintf("Unrecognized error code %v", code),
		}
	}
}
