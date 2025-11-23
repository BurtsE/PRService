package pr_service

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
)

func (s *Service) SetUserIsActive(ctx context.Context, user *model.User) error {
	exists, err := s.storage.UserExists(ctx, user.ID)
	if err != nil {
		s.logger.Errorf("Error checking if user exists: %v", err)
		return err
	}
	if !exists {
		return service.ErrResourceNotFound
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
	if user.IsActive == false {
		prs, err := s.GetReviewersPRs(ctx, user.ID)
		if err != nil {
			s.logger.Warnf("Error getting reviewers PRs: %v", err)
			return
		}

		for _, pr := range prs {
			_, err = s.ReassignPullRequestReviewer(ctx, pr.ID, user.ID)
			if err != nil {
				s.logger.Warnf("Error reassigning PR: %v", err)
			}
		}
	}
}
