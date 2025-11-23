package mocks

import (
	"PRService/internal/model"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreatePullRequest(ctx context.Context, request *model.PullRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *MockService) MergePullRequest(ctx context.Context, id model.PullRequestID) (*model.PullRequest, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PullRequest), args.Error(1)
}

func (m *MockService) ReassignPullRequestReviewer(ctx context.Context, pullRequestID model.PullRequestID, userID model.UserID) (model.PullRequest, error) {
	args := m.Called(ctx, pullRequestID, userID)
	// The first return value is a struct, not a pointer, so return an empty one on error.
	return args.Get(0).(model.PullRequest), args.Error(1)
}

func (m *MockService) GetReviewersPRs(ctx context.Context, userID model.UserID) ([]model.PullRequest, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PullRequest), args.Error(1)
}

func (m *MockService) CreateTeam(ctx context.Context, team *model.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockService) GetTeam(ctx context.Context, teamName model.TeamName) (*model.Team, error) {
	args := m.Called(ctx, teamName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Team), args.Error(1)
}

func (m *MockService) SetUserIsActive(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockService) GetStatistic(ctx context.Context) (model.Statistic, error) {
	args := m.Called(ctx)
	return args.Get(0).(model.Statistic), args.Error(1)
}

// The following methods are part of the service.Service interface but were missing from the mock.

func (m *MockService) TeamExists(ctx context.Context, teamName model.TeamName) (bool, error) {
	args := m.Called(ctx, teamName)
	return args.Bool(0), args.Error(1)
}

func (m *MockService) GetUser(ctx context.Context, id model.UserID) (model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.User), args.Error(1)
}
