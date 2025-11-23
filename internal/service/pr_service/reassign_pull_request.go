package prservice

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
	"math/rand"
	"slices"
	"time"
)

func (s *Service) ReassignPullRequestReviewer(ctx context.Context,
	pullRequestID model.PullRequestID, userID model.UserID) (model.PullRequest, error) {
	// Check if pr exists and not merged
	request, err := s.storage.GetPullRequest(ctx, pullRequestID)
	if err != nil {
		return model.PullRequest{}, err
	}

	if request.Status == model.PullRequestStatusMerged {
		return model.PullRequest{}, service.ErrPullRequestMerged
	}

	// Check if user exists
	user, err := s.storage.GetUser(ctx, userID)
	if err != nil {
		return model.PullRequest{}, err
	}

	// check if user is a reviewer
	var found bool
	for i := range request.Reviewers {
		if request.Reviewers[i] == userID {
			found = true
			break
		}
	}
	if !found {
		return model.PullRequest{}, service.ErrResourceNotFound
	}

	// Get user's team
	team, err := s.storage.GetTeam(ctx, user.TeamName)
	if err != nil {
		return model.PullRequest{}, err
	}

	newReviewer := s.getAvailableReveiwer(*team, request)
	if newReviewer.ID == "" {
		return model.PullRequest{}, service.ErrReviewersUnavailable
	}

	err = s.storage.ReassignPullRequestReviewer(ctx, &request, userID, newReviewer.ID)
	if err != nil {
		return model.PullRequest{}, err
	}
	

	return request, nil
}

func (s *Service) getAvailableReveiwer(team model.Team, request model.PullRequest) model.User {
	var availableReveiwers []model.User

	for _, member := range team.Members {
		if !member.IsActive {
			continue
		}

		if member.ID == request.AuthorID {
			continue
		}

		if slices.ContainsFunc(request.Reviewers, func(userID model.UserID) bool {
			return userID == member.ID
		}) {
			continue
		}

		availableReveiwers = append(availableReveiwers, member)
	}

	if len(availableReveiwers) == 0 {
		return model.User{}
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	index := r.Intn(len(availableReveiwers))

	return availableReveiwers[index]
}
