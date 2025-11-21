package service

import (
	"PRService/internal/model"
	"context"
)

type Service interface {
	CreateTeam(context.Context, *model.Team) error
	GetTeam(teamName string) (*model.Team, error)

	SetUserIsActive(model.UserID) error
	GetReviewersPRs(model.UserID) ([]model.PullRequest, error)

	CreatePullRequest(model.PullRequest) error
	MergePullRequest(model.PullRequestID) (*model.PullRequest, error)
	ReassignPullRequestReviewer(model.PullRequestID, model.UserID) error
}
