package router

import (
	"PRService/internal/mocks"
	"PRService/internal/model"
	"PRService/internal/service"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTest(t *testing.T) (*Router, *mocks.MockService) {
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	mockService := new(mocks.MockService)
	router := NewRouter(nil, logger, mockService)
	router.SetupRoutes()
	return router, mockService
}

func TestCreateTeam(t *testing.T) {
	router, mockService := setupTest(t)

	team := model.Team{
		Name: "test-team",
		Members: []model.User{
			{ID: "user-1", Name: "Test User"},
		},
	}
	body, _ := json.Marshal(team)

	t.Run("success", func(t *testing.T) {
		mockService.On("CreateTeam", mock.Anything, &team).Return(nil).Once()

		req, _ := http.NewRequest("POST", "/team/add", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("team exists", func(t *testing.T) {
		mockService.On("CreateTeam", mock.Anything, &team).Return(service.ErrTeamExists).Once()

		req, _ := http.NewRequest("POST", "/team/add", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/team/add", bytes.NewReader([]byte(`{"invalid`)))
		req.Header.Set("Content-Type", "application/json")

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode) // Fiber's default for bind error
	})

	t.Run("invalid team model", func(t *testing.T) {
		invalidTeam := model.Team{Name: ""} // Invalid because name is empty
		body, _ := json.Marshal(invalidTeam)
		req, _ := http.NewRequest("POST", "/team/add", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode) // Mapped to ResourceNotFound
	})
}

func TestGetTeam(t *testing.T) {
	router, mockService := setupTest(t)

	teamName := model.TeamName("test-team")
	team := &model.Team{Name: teamName}

	t.Run("success", func(t *testing.T) {
		mockService.On("GetTeam", mock.Anything, teamName).Return(team, nil).Once()

		req, _ := http.NewRequest("GET", "/team/get?team_name=test-team", nil)

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var resultTeam model.Team
		bodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &resultTeam)
		assert.Equal(t, team.Name, resultTeam.Name)

		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("GetTeam", mock.Anything, teamName).Return(nil, service.ErrResourceNotFound).Once()

		req, _ := http.NewRequest("GET", "/team/get?team_name=test-team", nil)

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("missing query param", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/team/get", nil)

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestSetUserIsActive(t *testing.T) {
	router, mockService := setupTest(t)

	user := model.User{ID: "user-1", IsActive: false}
	body, _ := json.Marshal(user)

	t.Run("success", func(t *testing.T) {
		mockService.On("SetUserIsActive", mock.Anything, &user).Return(nil).Once()

		req, _ := http.NewRequest("POST", "/users/setIsActive", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockService.On("SetUserIsActive", mock.Anything, &user).Return(service.ErrResourceNotFound).Once()

		req, _ := http.NewRequest("POST", "/users/setIsActive", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}

func TestCreatePullRequest(t *testing.T) {
	router, mockService := setupTest(t)

	pr := model.PullRequest{
		ID:       "pr-1",
		Name:     "feat: new feature",
		AuthorID: "user-1",
	}
	body, _ := json.Marshal(pr)

	t.Run("success", func(t *testing.T) {
		mockService.On("CreatePullRequest", mock.Anything, &pr).Return(nil).Once()

		req, _ := http.NewRequest("POST", "/pullRequest/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("pr exists", func(t *testing.T) {
		mockService.On("CreatePullRequest", mock.Anything, &pr).Return(service.ErrPullRequestExists).Once()

		req, _ := http.NewRequest("POST", "/pullRequest/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}

func TestMergePullRequest(t *testing.T) {
	router, mockService := setupTest(t)

	reqBody := map[string]string{"pull_request_id": "pr-1"}
	body, _ := json.Marshal(reqBody)
	prID := model.PullRequestID("pr-1")
	mergedPR := &model.PullRequest{ID: prID, Status: model.PullRequestStatusMerged}

	t.Run("success", func(t *testing.T) {
		mockService.On("MergePullRequest", mock.Anything, prID).Return(mergedPR, nil).Once()

		req, _ := http.NewRequest("POST", "/pullRequest/merge", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var resultPR model.PullRequest
		bodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &resultPR)
		assert.Equal(t, mergedPR.Status, resultPR.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("MergePullRequest", mock.Anything, prID).Return(nil, service.ErrResourceNotFound).Once()

		req, _ := http.NewRequest("POST", "/pullRequest/merge", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}

func TestReassignPullRequestReviewer(t *testing.T) {
	router, mockService := setupTest(t)

	reqBody := map[string]string{"pull_request_id": "pr-1", "old_reviewer_id": "user-2"}
	body, _ := json.Marshal(reqBody)
	prID := model.PullRequestID("pr-1")
	userID := model.UserID("user-2")
	reassignedPR := model.PullRequest{ID: prID, Reviewers: []model.UserID{"user-3"}}

	t.Run("success", func(t *testing.T) {
		mockService.On("ReassignPullRequestReviewer", mock.Anything, prID, userID).Return(reassignedPR, nil).Once()

		req, _ := http.NewRequest("POST", "/pullRequest/reassign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var resultPR model.PullRequest
		bodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &resultPR)
		assert.Contains(t, resultPR.Reviewers, model.UserID("user-3"))

		mockService.AssertExpectations(t)
	})

	t.Run("pr merged", func(t *testing.T) {
		mockService.On("ReassignPullRequestReviewer", mock.Anything, prID, userID).Return(model.PullRequest{}, service.ErrPullRequestMerged).Once()

		req, _ := http.NewRequest("POST", "/pullRequest/reassign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}

func TestGetReview(t *testing.T) {
	router, mockService := setupTest(t)

	userID := model.UserID("user-1")
	prs := []model.PullRequest{
		{ID: "pr-1", Name: "PR 1", AuthorID: "user-2", Status: model.PullRequestStatusOpen},
	}

	t.Run("success", func(t *testing.T) {
		mockService.On("GetReviewersPRs", mock.Anything, userID).Return(prs, nil).Once()

		req, _ := http.NewRequest("GET", "/users/getReview?user_id=user-1", nil)

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result struct {
			UserID       string           `json:"user_id"`
			PullRequests []PullRequestDto `json:"pullRequests"`
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(bodyBytes, &result)

		assert.Equal(t, string(userID), result.UserID)
		assert.Len(t, result.PullRequests, 1)
		assert.Equal(t, prs[0].ID, result.PullRequests[0].PullRequestID)

		mockService.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockService.On("GetReviewersPRs", mock.Anything, userID).Return(nil, service.ErrResourceNotFound).Once()

		req, _ := http.NewRequest("GET", "/users/getReview?user_id=user-1", nil)

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("missing user_id", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/getReview", nil)

		resp, err := router.app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestPing(t *testing.T) {
	router, _ := setupTest(t)

	req, _ := http.NewRequest("GET", "/ping", nil)

	resp, err := router.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `"pong"`, string(bodyBytes))
}
