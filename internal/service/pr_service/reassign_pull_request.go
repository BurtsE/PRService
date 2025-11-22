package pr_service

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
)

func (s *Service) ReassignPullRequestReviewer(ctx context.Context,
	pullRequestID model.PullRequestID, userID model.UserID) (model.PullRequest, error) {

	request, err := s.storage.GetPullRequest(ctx, pullRequestID)
	if err != nil {
		return model.PullRequest{}, err
	}

	if request.Status == model.PullRequestStatusMerged {
		return model.PullRequest{}, service.ErrPullRequestMerged
	}

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
	err = s.storage.ReassignPullRequestReviewer(ctx, &request, userID)
	if err != nil {
		return model.PullRequest{}, err
	}

	return request, nil
}
