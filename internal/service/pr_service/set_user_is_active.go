package prservice

import (
	"PRService/internal/model"
	"context"
)

// SetUserIsActive changes user's is_active flag, then tries to reassign reviewed pull requests
func (s *Service) SetUserIsActive(ctx context.Context, user *model.User) error {
	_, err := s.storage.GetUser(ctx, user.ID)
	if err != nil {
		s.logger.Errorf("SetUserIsActive: could not get user: %v", err)
		return err
	}

	err = s.storage.SetUserIsActive(ctx, user)
	if err != nil {

		s.logger.Errorf("SetUserIsActive: could not set user is_active: %v", err)
		return err
	}

	if !user.IsActive {
		s.reassignInactiveUsersPrs(ctx, user)
	}

	return nil
}

func (s *Service) reassignInactiveUsersPrs(ctx context.Context, user *model.User) {
	if !user.IsActive {
		prs, err := s.GetReviewersPRs(ctx, user.ID)
		if err != nil {
			return
		}

		for _, pr := range prs {
			if pr.Status == model.PullRequestStatusMerged {
				continue
			}
			_, _ = s.ReassignPullRequestReviewer(ctx, pr.ID, user.ID)
		}
	}
}
