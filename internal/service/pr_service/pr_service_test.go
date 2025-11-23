package prservice

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePullRequest(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)

	pr := &model.PullRequest{ID: "pr-1", AuthorID: "user-1"}

	t.Run("success", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mockStorage.On("GetPullRequest", ctx, pr.ID).Return(model.PullRequest{}, service.ErrResourceNotFound).Once()
		mockStorage.On("GetUser", ctx, pr.AuthorID).Return(model.User{ID: pr.AuthorID}, nil).Once()
		mockStorage.On("CreatePullRequest", ctx, pr).Return(nil).Once()

		err := s.CreatePullRequest(ctx, pr)

		assert.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("pull request already exists", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mockStorage.On("GetPullRequest", ctx, pr.ID).Return(*pr, nil).Once()

		err := s.CreatePullRequest(ctx, pr)

		assert.Equal(t, service.ErrPullRequestExists, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("author does not exist", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mockStorage.On("GetPullRequest", ctx, pr.ID).Return(model.PullRequest{}, service.ErrResourceNotFound).Once()
		mockStorage.On("GetUser", ctx, pr.AuthorID).Return(model.User{}, service.ErrResourceNotFound).Once()

		err := s.CreatePullRequest(ctx, pr)

		assert.Equal(t, service.ErrResourceNotFound, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("storage error on create", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		storageErr := errors.New("storage error")
		mockStorage.On("GetPullRequest", ctx, pr.ID).Return(model.PullRequest{}, service.ErrResourceNotFound).Once()
		mockStorage.On("GetUser", ctx, pr.AuthorID).Return(model.User{ID: pr.AuthorID}, nil).Once()
		mockStorage.On("CreatePullRequest", ctx, pr).Return(storageErr).Once()

		err := s.CreatePullRequest(ctx, pr)

		assert.Equal(t, storageErr, err)
		mockStorage.AssertExpectations(t)
	})
}

func TestCreateTeam(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)

	team := &model.Team{Name: "team-a"}

	t.Run("success", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mockStorage.On("TeamExists", ctx, team.Name).Return(false, nil).Once()
		mockStorage.On("CreateTeam", ctx, team).Return(nil).Once()

		err := s.CreateTeam(ctx, team)

		assert.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("team already exists", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mockStorage.On("TeamExists", ctx, team.Name).Return(true, nil).Once()

		err := s.CreateTeam(ctx, team)

		assert.Equal(t, service.ErrTeamExists, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("storage error", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		storageErr := errors.New("storage error")
		mockStorage.On("TeamExists", ctx, team.Name).Return(false, nil).Once()
		mockStorage.On("CreateTeam", ctx, team).Return(storageErr).Once()

		err := s.CreateTeam(ctx, team)

		assert.Equal(t, storageErr, err)
		mockStorage.AssertExpectations(t)
	})
}

func TestGetTeam(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)

	teamName := model.TeamName("team-a")
	expectedTeam := &model.Team{Name: teamName}

	t.Run("success", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mockStorage.On("TeamExists", ctx, teamName).Return(true, nil).Once()
		mockStorage.On("GetTeam", ctx, teamName).Return(expectedTeam, nil).Once()

		team, err := s.GetTeam(ctx, teamName)

		assert.NoError(t, err)
		assert.Equal(t, expectedTeam, team)
		mockStorage.AssertExpectations(t)
	})

	t.Run("team not found", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mockStorage.On("TeamExists", ctx, teamName).Return(false, nil).Once()

		team, err := s.GetTeam(ctx, teamName)

		assert.Equal(t, service.ErrResourceNotFound, err)
		assert.Nil(t, team)
		mockStorage.AssertExpectations(t)
	})
}

func TestMergePullRequest(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)

	prID := model.PullRequestID("pr-1")
	pr := model.PullRequest{ID: prID}
	mergedPR := &model.PullRequest{ID: prID, Status: model.PullRequestStatusMerged}

	t.Run("success", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mockStorage.On("GetPullRequest", ctx, prID).Return(pr, nil).Once()
		mockStorage.On("MergePullRequest", ctx, prID).Return(mergedPR, nil).Once()

		result, err := s.MergePullRequest(ctx, prID)

		assert.NoError(t, err)
		assert.Equal(t, mergedPR, result)
		mockStorage.AssertExpectations(t)
	})

	t.Run("pr not found", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mockStorage.On("GetPullRequest", ctx, prID).Return(model.PullRequest{}, service.ErrResourceNotFound).Once()

		result, err := s.MergePullRequest(ctx, prID)

		assert.Equal(t, service.ErrResourceNotFound, err)
		assert.Nil(t, result)
		mockStorage.AssertExpectations(t)
	})
}

func TestReassignPullRequestReviewer(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)

	prID := model.PullRequestID("pr-1")
	oldReviewerID := model.UserID("user-1")
	newReviewerID := model.UserID("user-3")
	pr := model.PullRequest{ID: prID, AuthorID: "author-1", Status: model.PullRequestStatusOpen, Reviewers: []model.UserID{oldReviewerID, "user-2"}}
	user := model.User{ID: oldReviewerID, TeamName: "team-a"}
	team := &model.Team{Name: "team-a", Members: []model.User{
		{ID: "author-1"},
		{ID: oldReviewerID, IsActive: true},
		{ID: "user-2", IsActive: true},
		{ID: newReviewerID, IsActive: true},
	}}

	t.Run("success", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mockStorage.On("GetPullRequest", ctx, prID).Return(pr, nil).Once()
		mockStorage.On("GetUser", ctx, oldReviewerID).Return(user, nil).Once()
		mockStorage.On("GetTeam", ctx, user.TeamName).Return(team, nil).Once()
		mockStorage.On("ReassignPullRequestReviewer", ctx, &pr, oldReviewerID, newReviewerID).Return(nil).Once()

		result, err := s.ReassignPullRequestReviewer(ctx, prID, oldReviewerID)

		assert.NoError(t, err)
		assert.Contains(t, result.Reviewers, newReviewerID)
		assert.NotContains(t, result.Reviewers, oldReviewerID)
		mockStorage.AssertExpectations(t)
	})

	t.Run("pr not found", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mockStorage.On("GetPullRequest", ctx, prID).Return(model.PullRequest{}, service.ErrResourceNotFound).Once()

		_, err := s.ReassignPullRequestReviewer(ctx, prID, oldReviewerID)

		assert.Equal(t, service.ErrResourceNotFound, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("pr already merged", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mergedPR := model.PullRequest{ID: prID, Status: model.PullRequestStatusMerged}
		mockStorage.On("GetPullRequest", ctx, prID).Return(mergedPR, nil).Once()

		_, err := s.ReassignPullRequestReviewer(ctx, prID, oldReviewerID)

		assert.Equal(t, service.ErrPullRequestMerged, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("reviewer not found on pr", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		otherUser := model.UserID("user-99")
		mockStorage.On("GetPullRequest", ctx, prID).Return(pr, nil).Once()
		mockStorage.On("GetUser", ctx, otherUser).Return(model.User{ID: otherUser}, nil).Once()

		_, err := s.ReassignPullRequestReviewer(ctx, prID, otherUser)

		assert.Equal(t, service.ErrResourceNotFound, err)
		mockStorage.AssertExpectations(t)
	})
}

func TestSetUserIsActive(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)

	teamName := model.TeamName("team-a")
	user := &model.User{ID: "user-1", IsActive: false, TeamName: teamName}

	t.Run("success - set inactive", func(t *testing.T) {
		mockStorage, s := newTestService(logger)

		pr := model.PullRequest{ID: "pr-1", AuthorID: "author-1", Status: model.PullRequestStatusOpen, Reviewers: []model.UserID{user.ID}}
		team := &model.Team{Name: teamName, Members: []model.User{{ID: "new-reviewer-1", IsActive: true}}}

		mockStorage.On("GetUser", ctx, user.ID).Return(*user, nil)
		mockStorage.On("SetUserIsActive", ctx, user).Return(nil)
		// For reassignInactiveUsersPrs
		mockStorage.On("GetReviewersPRs", ctx, user.ID).Return([]model.PullRequest{pr}, nil)
		// Mocks for the inner ReassignPullRequestReviewer call
		mockStorage.On("GetPullRequest", ctx, pr.ID).Return(pr, nil)
		mockStorage.On("GetUser", ctx, user.ID).Return(model.User{ID: user.ID, TeamName: "team-a"}, nil)
		mockStorage.On("GetTeam", ctx, model.TeamName("team-a")).Return(team, nil)
		mockStorage.On("ReassignPullRequestReviewer", ctx, mock.Anything, user.ID, mock.AnythingOfType("model.UserID")).Return(nil)

		err := s.SetUserIsActive(ctx, user)

		assert.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("success - set active", func(t *testing.T) {
		mockStorage, s := newTestService(logger)

		activeUser := &model.User{ID: "user-1", IsActive: true}
		mockStorage.On("GetUser", ctx, activeUser.ID).Return(*activeUser, nil)
		mockStorage.On("SetUserIsActive", ctx, activeUser).Return(nil)

		err := s.SetUserIsActive(ctx, activeUser)

		assert.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockStorage, s := newTestService(logger)

		mockStorage.On("GetUser", ctx, user.ID).Return(model.User{}, service.ErrResourceNotFound).Once()

		err := s.SetUserIsActive(ctx, user)

		assert.Equal(t, service.ErrResourceNotFound, err)
		mockStorage.AssertExpectations(t)
	})
}

func TestGetReviewersPRs(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)

	userID := model.UserID("user-1")
	expectedPRs := []model.PullRequest{{ID: "pr-1"}, {ID: "pr-2"}}

	t.Run("success", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mockStorage.On("GetUser", ctx, userID).Return(model.User{ID: userID}, nil).Once()
		mockStorage.On("GetReviewersPRs", ctx, userID).Return(expectedPRs, nil).Once()

		prs, err := s.GetReviewersPRs(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPRs, prs)
		mockStorage.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockStorage, s := newTestService(logger)
		mockStorage.On("GetUser", ctx, userID).Return(model.User{}, service.ErrResourceNotFound).Once()

		prs, err := s.GetReviewersPRs(ctx, userID)

		assert.Equal(t, service.ErrResourceNotFound, err)
		assert.Nil(t, prs)
		mockStorage.AssertExpectations(t)
	})
}

func newTestService(logger *logrus.Logger) (*MockStorage, *Service) {
	mockStorage := new(MockStorage)
	return mockStorage, NewService(logger, mockStorage)
}
