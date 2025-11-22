package storage

import (
	"PRService/internal/model"
	"context"
)

type Storage interface {
	TeamExists(context.Context, model.TeamName) (bool, error)
	CreateTeam(context.Context, *model.Team) error
	GetTeam(context.Context, model.TeamName) (*model.Team, error)

	UserExists(context.Context, model.UserID) (bool, error)
	SetUserIsActive(context.Context, *model.User) error
	GetReviewersPRs(context.Context, model.UserID) ([]model.PullRequest, error)

	PullRequestExists(context.Context, model.PullRequestID) (bool, error)
	CreatePullRequest(context.Context, model.PullRequest) error
	MergePullRequest(context.Context, model.PullRequestID) (*model.PullRequest, error)
	ReassignPullRequestReviewer(context.Context, model.PullRequestID, model.UserID) error
}
