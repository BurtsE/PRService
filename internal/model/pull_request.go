package model

import "time"

type PullRequestStatus string
type PullRequestID string

const (
	PullRequestStatusOpen   PullRequestStatus = "OPEN"
	PullRequestStatusMerged PullRequestStatus = "MERGED"
)

type PullRequest struct {
	ID        PullRequestID     `json:"pull_request_id"`
	Name      string            `json:"pull_request_name"`
	AuthorID  UserID            `json:"author_id"`
	Reviewers []UserID          `json:"assigned_reviewers"`
	Status    PullRequestStatus `json:"status"`
	CreatedAt time.Time         `json:"-"`
	MergedAt  time.Time         `json:"-"`
}

func (p *PullRequest) Valid() bool {
	return len(p.ID) != 0 && len(p.AuthorID) != 0 && len(p.Name) != 0
}

func (p *PullRequest) Init() {
	p.CreatedAt = time.Now()
	p.Status = PullRequestStatusOpen
	p.Reviewers = make([]UserID, 0)
}
