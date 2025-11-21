package pr_service

import (
	"PRService/internal/model"
	"PRService/internal/service"
)

var _ service.Service = (*Service)(nil)

type Service struct {
}

func (s Service) CreateTeam(team *model.Team) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetTeam(teamName string) (*model.Team, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) SetUserIsActive(id model.UserID) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetReviewersPRs(id model.UserID) ([]model.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) CreatePullRequest(request model.PullRequest) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) MergePullRequest(id model.PullRequestID) (*model.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) ReassignPullRequestReviewer(id model.PullRequestID, id2 model.UserID) error {
	//TODO implement me
	panic("implement me")
}
