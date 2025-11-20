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
	AuthorId  UserId
	Reviewers []UserId
	Status    PullRequestStatus
	CreatedAt time.Time
	MergedAt  time.Time
}
