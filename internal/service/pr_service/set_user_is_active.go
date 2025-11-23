package prservice

import (
	"PRService/internal/model"
	"context"
)

// SetUserIsActive changes user's is_active flag, then tries to reassign reviewed pull requests
func (s *Service) SetUserIsActive(ctx context.Context, user *model.User) error {
	_, err := s.storage.GetUser(ctx, user.ID)
	if err != nil {
		return err
	}

	err = s.storage.SetUserIsActive(ctx, user)
	if err != nil {
		s.logger.Errorf("Error setting user isActive: %v", err)
		return err
	}

	s.reassignInactiveUsersPrs(ctx, user)

	return nil
}

func (s *Service) reassignInactiveUsersPrs(ctx context.Context, user *model.User) {
	if !user.IsActive {
		prs, err := s.GetReviewersPRs(ctx, user.ID)
		if err != nil {
			s.logger.Warnf("Error getting reviewers PRs: %v", err)
			return
		}

		for _, pr := range prs {
			if pr.Status == model.PullRequestStatusMerged {
				continue
			}

			_, err = s.ReassignPullRequestReviewer(ctx, pr.ID, user.ID)
			if err != nil {
				s.logger.Warnf("Error reassigning PR: %v", err)
			}
		}
	}
}
