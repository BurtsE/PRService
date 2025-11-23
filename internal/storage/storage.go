package storage

import (
	"PRService/internal/model"
	"context"
)

type Storage interface {
	TeamExists(context.Context, model.TeamName) (bool, error)
	CreateTeam(context.Context, *model.Team) error
	GetTeam(context.Context, model.TeamName) (*model.Team, error)

	GetUser(context.Context, model.UserID) (model.User, error)
	SetUserIsActive(context.Context, *model.User) error
	GetReviewersPRs(context.Context, model.UserID) ([]model.PullRequest, error)

	GetPullRequest(context.Context, model.PullRequestID) (model.PullRequest, error)
	CreatePullRequest(context.Context, *model.PullRequest) error
	MergePullRequest(context.Context, model.PullRequestID) (*model.PullRequest, error)
	ReassignPullRequestReviewer(context.Context, *model.PullRequest, model.UserID, model.UserID) error
}
