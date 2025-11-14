package domain

import "time"

type PRStatus string

const (
	PROpen   PRStatus = "OPEN"
	PRMerged PRStatus = "MERGED"
)

type PullRequest struct {
	ID        string    `json:"pull_request_id" db:"id"`
	Name      string    `json:"pull_request_name" db:"name"`
	AuthorID  string    `json:"author_id" db:"author_id"`
	Status    PRStatus  `json:"status" db:"status"`
	Reviewers []string  `json:"assigned_reviewers"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	MergedAt  time.Time `json:"mergedAt" db:"merged_at"`
}

type PullRequestShort struct {
	ID       string   `json:"pull_request_id" db:"id"`
	Name     string   `json:"pull_request_name" db:"name"`
	AuthorID string   `json:"author_id" db:"author_id"`
	Status   PRStatus `json:"status" db:"status"`
}
