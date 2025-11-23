package service

import (
	"PRService/internal/model"
	"context"
	"fmt"
)

var (
	ErrTeamExists           = fmt.Errorf("team already exists")
	ErrPullRequestExists    = fmt.Errorf("pull request already exists")
	ErrPullRequestMerged    = fmt.Errorf("cannot reassign on merged PR")
	ErrReviewerNotAssigned  = fmt.Errorf("reviewer not assigned")
	ErrReviewersUnavailable = fmt.Errorf("reviewers are not available")
	ErrResourceNotFound     = fmt.Errorf("resource not found")
)

type Service interface {
	CreateTeam(context.Context, *model.Team) error
	GetTeam(context.Context, model.TeamName) (*model.Team, error)

	SetUserIsActive(context.Context, *model.User) error
	GetReviewersPRs(context.Context, model.UserID) ([]model.PullRequest, error)

	CreatePullRequest(context.Context, *model.PullRequest) error
	MergePullRequest(context.Context, model.PullRequestID) (*model.PullRequest, error)
	ReassignPullRequestReviewer(context.Context, model.PullRequestID, model.UserID) (model.PullRequest, error)

	GetStatistic(context.Context) (model.Statistic, error)
}
