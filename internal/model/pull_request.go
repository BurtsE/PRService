package model

import "time"

type PullRequestStatus string
type PullRequestID string

const (
	PullRequestStatusOpen   PullRequestStatus = "OPEN"
	PullRequestStatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	Id        PullRequestID
	AuthorId  UserID
	Reviewers []UserID
	Status    PullRequestStatus
	CreatedAt time.Time
	MergedAt  time.Time
}

func (p *PullRequest) Valid() bool {
	return p.Id != "" && p.AuthorId != ""
}
