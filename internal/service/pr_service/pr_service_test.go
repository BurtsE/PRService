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
	mockStorage := new(MockStorage)
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	s := NewService(logger, mockStorage)

	pr := &model.PullRequest{ID: "pr-1", AuthorID: "user-1"}

	t.Run("success", func(t *testing.T) {
		mockStorage.On("GetPullRequest", ctx, pr.ID).Return(model.PullRequest{}, service.ErrResourceNotFound).Once()
		mockStorage.On("UserExists", ctx, pr.AuthorID).Return(true, nil).Once()
		mockStorage.On("CreatePullRequest", ctx, pr).Return(nil).Once()

		err := s.CreatePullRequest(ctx, pr)

		assert.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("pull request already exists", func(t *testing.T) {
		mockStorage.On("GetPullRequest", ctx, pr.ID).Return(*pr, nil).Once()

		err := s.CreatePullRequest(ctx, pr)

		assert.Equal(t, service.ErrPullRequestExists, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("author does not exist", func(t *testing.T) {
		mockStorage.On("GetPullRequest", ctx, pr.ID).Return(model.PullRequest{}, service.ErrResourceNotFound).Once()
		mockStorage.On("UserExists", ctx, pr.AuthorID).Return(false, nil).Once()

		err := s.CreatePullRequest(ctx, pr)

		assert.Equal(t, service.ErrResourceNotFound, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("storage error on create", func(t *testing.T) {
		storageErr := errors.New("storage error")
		mockStorage.On("GetPullRequest", ctx, pr.ID).Return(model.PullRequest{}, service.ErrResourceNotFound).Once()
		mockStorage.On("UserExists", ctx, pr.AuthorID).Return(true, nil).Once()
		mockStorage.On("CreatePullRequest", ctx, pr).Return(storageErr).Once()

		err := s.CreatePullRequest(ctx, pr)

		assert.Equal(t, storageErr, err)
		mockStorage.AssertExpectations(t)
	})
}

func TestCreateTeam(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	s := NewService(logger, mockStorage)

	team := &model.Team{Name: "team-a"}

	t.Run("success", func(t *testing.T) {
		mockStorage.On("TeamExists", ctx, team.Name).Return(false, nil).Once()
		mockStorage.On("CreateTeam", ctx, team).Return(nil).Once()

		err := s.CreateTeam(ctx, team)

		assert.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("team already exists", func(t *testing.T) {
		mockStorage.On("TeamExists", ctx, team.Name).Return(true, nil).Once()

		err := s.CreateTeam(ctx, team)

		assert.Equal(t, service.ErrTeamExists, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("storage error", func(t *testing.T) {
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
	mockStorage := new(MockStorage)
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	s := NewService(logger, mockStorage)

	teamName := model.TeamName("team-a")
	expectedTeam := &model.Team{Name: teamName}

	t.Run("success", func(t *testing.T) {
		mockStorage.On("TeamExists", ctx, teamName).Return(true, nil).Once()
		mockStorage.On("GetTeam", ctx, teamName).Return(expectedTeam, nil).Once()

		team, err := s.GetTeam(ctx, teamName)

		assert.NoError(t, err)
		assert.Equal(t, expectedTeam, team)
		mockStorage.AssertExpectations(t)
	})

	t.Run("team not found", func(t *testing.T) {
		mockStorage.On("TeamExists", ctx, teamName).Return(false, nil).Once()

		team, err := s.GetTeam(ctx, teamName)

		assert.Equal(t, service.ErrResourceNotFound, err)
		assert.Nil(t, team)
		mockStorage.AssertExpectations(t)
	})
}

func TestMergePullRequest(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	s := NewService(logger, mockStorage)

	prID := model.PullRequestID("pr-1")
	pr := model.PullRequest{ID: prID}
	mergedPR := &model.PullRequest{ID: prID, Status: model.PullRequestStatusMerged}

	t.Run("success", func(t *testing.T) {
		mockStorage.On("GetPullRequest", ctx, prID).Return(pr, nil).Once()
		mockStorage.On("MergePullRequest", ctx, prID).Return(mergedPR, nil).Once()

		result, err := s.MergePullRequest(ctx, prID)

		assert.NoError(t, err)
		assert.Equal(t, mergedPR, result)
		mockStorage.AssertExpectations(t)
	})

	t.Run("pr not found", func(t *testing.T) {
		mockStorage.On("GetPullRequest", ctx, prID).Return(model.PullRequest{}, service.ErrResourceNotFound).Once()

		result, err := s.MergePullRequest(ctx, prID)

		assert.Equal(t, service.ErrResourceNotFound, err)
		assert.Nil(t, result)
		mockStorage.AssertExpectations(t)
	})
}

func TestReassignPullRequestReviewer(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	s := NewService(logger, mockStorage)

	prID := model.PullRequestID("pr-1")
	userID := model.UserID("user-1")
	pr := model.PullRequest{ID: prID, Status: model.PullRequestStatusOpen, Reviewers: []model.UserID{userID, "user-2"}}

	t.Run("success", func(t *testing.T) {
		mockStorage.On("GetPullRequest", ctx, prID).Return(pr, nil).Once()
		mockStorage.On("ReassignPullRequestReviewer", ctx, &pr, userID).Return(nil).Once()

		result, err := s.ReassignPullRequestReviewer(ctx, prID, userID)

		assert.NoError(t, err)
		assert.Equal(t, pr, result)
		mockStorage.AssertExpectations(t)
	})

	t.Run("pr not found", func(t *testing.T) {
		mockStorage.On("GetPullRequest", ctx, prID).Return(model.PullRequest{}, service.ErrResourceNotFound).Once()

		_, err := s.ReassignPullRequestReviewer(ctx, prID, userID)

		assert.Equal(t, service.ErrResourceNotFound, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("pr already merged", func(t *testing.T) {
		mergedPR := model.PullRequest{ID: prID, Status: model.PullRequestStatusMerged}
		mockStorage.On("GetPullRequest", ctx, prID).Return(mergedPR, nil).Once()

		_, err := s.ReassignPullRequestReviewer(ctx, prID, userID)

		assert.Equal(t, service.ErrPullRequestMerged, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("reviewer not found on pr", func(t *testing.T) {
		otherUser := model.UserID("user-99")
		mockStorage.On("GetPullRequest", ctx, prID).Return(pr, nil).Once()

		_, err := s.ReassignPullRequestReviewer(ctx, prID, otherUser)

		assert.Equal(t, service.ErrResourceNotFound, err)
		mockStorage.AssertExpectations(t)
	})
}

func TestSetUserIsActive(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)

	user := &model.User{ID: "user-1", IsActive: false}

	t.Run("success - set inactive", func(t *testing.T) {
		mockStorage := new(MockStorage)
		s := NewService(logger, mockStorage)

		pr := model.PullRequest{ID: "pr-1", Status: model.PullRequestStatusOpen, Reviewers: []model.UserID{user.ID}}
		mockStorage.On("UserExists", ctx, user.ID).Return(true, nil)
		mockStorage.On("SetUserIsActive", ctx, user).Return(nil)
		// For reassignInactiveUsersPrs
		mockStorage.On("GetReviewersPRs", ctx, user.ID).Return([]model.PullRequest{pr}, nil)
		mockStorage.On("GetPullRequest", ctx, pr.ID).Return(pr, nil)
		mockStorage.On("ReassignPullRequestReviewer", ctx, mock.Anything, user.ID).Return(nil)

		err := s.SetUserIsActive(ctx, user)

		assert.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("success - set active", func(t *testing.T) {
		mockStorage := new(MockStorage)
		s := NewService(logger, mockStorage)

		activeUser := &model.User{ID: "user-1", IsActive: true}
		mockStorage.On("UserExists", ctx, activeUser.ID).Return(true, nil)
		mockStorage.On("SetUserIsActive", ctx, activeUser).Return(nil)

		err := s.SetUserIsActive(ctx, activeUser)

		assert.NoError(t, err)
		mockStorage.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockStorage := new(MockStorage)
		s := NewService(logger, mockStorage)

		mockStorage.On("UserExists", ctx, user.ID).Return(false, nil).Once()

		err := s.SetUserIsActive(ctx, user)

		assert.Equal(t, service.ErrResourceNotFound, err)
		mockStorage.AssertExpectations(t)
	})
}

func TestGetReviewersPRs(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	s := NewService(logger, mockStorage)

	userID := model.UserID("user-1")
	expectedPRs := []model.PullRequest{{ID: "pr-1"}, {ID: "pr-2"}}

	t.Run("success", func(t *testing.T) {
		mockStorage.On("UserExists", ctx, userID).Return(true, nil).Once()
		mockStorage.On("GetReviewersPRs", ctx, userID).Return(expectedPRs, nil).Once()

		prs, err := s.GetReviewersPRs(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPRs, prs)
		mockStorage.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockStorage.On("UserExists", ctx, userID).Return(false, nil).Once()

		prs, err := s.GetReviewersPRs(ctx, userID)

		assert.Equal(t, service.ErrResourceNotFound, err)
		assert.Nil(t, prs)
		mockStorage.AssertExpectations(t)
	})
}
