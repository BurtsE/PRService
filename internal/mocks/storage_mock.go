package mocks

import (
	"PRService/internal/model"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetPullRequest(ctx context.Context, pullRequestID model.PullRequestID) (model.PullRequest, error) {
	args := m.Called(ctx, pullRequestID)
	return args.Get(0).(model.PullRequest), args.Error(1)
}

func (m *MockStorage) CreatePullRequest(ctx context.Context, request *model.PullRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *MockStorage) MergePullRequest(ctx context.Context, requestID model.PullRequestID) (*model.PullRequest, error) {
	args := m.Called(ctx, requestID)
	// Handle nil for the first return value if an error occurs
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.PullRequest), args.Error(1)
}

func (m *MockStorage) ReassignPullRequestReviewer(ctx context.Context, pullRequest *model.PullRequest, oldReviewerID, newReviewerID model.UserID) error {
	args := m.Called(ctx, pullRequest, oldReviewerID, newReviewerID)
	return args.Error(0)
}

func (m *MockStorage) GetReviewersPRs(ctx context.Context, userID model.UserID) ([]model.PullRequest, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.PullRequest), args.Error(1)
}

func (m *MockStorage) CreateTeam(ctx context.Context, team *model.Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockStorage) GetTeam(ctx context.Context, teamName model.TeamName) (*model.Team, error) {
	args := m.Called(ctx, teamName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Team), args.Error(1)
}

func (m *MockStorage) TeamExists(ctx context.Context, teamName model.TeamName) (bool, error) {
	args := m.Called(ctx, teamName)
	return args.Bool(0), args.Error(1)
}

func (m *MockStorage) SetUserIsActive(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockStorage) GetUser(ctx context.Context, id model.UserID) (model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.User), args.Error(1)
}
