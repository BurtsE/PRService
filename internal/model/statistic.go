package model

import "time"

// Statistic holds aggregated data about the service.
type Statistic struct {
	TotalPRs    int `json:"total_prs"`
	OpenPRs     int `json:"open_prs"`
	MergedPRs   int `json:"merged_prs"`
	TotalUsers  int `json:"total_users"`
	ActiveUsers int `json:"active_users"`
	TotalTeams  int `json:"total_teams"`
	// AvgMergeTimeSeconds is the average time in seconds from PR creation to merge.
	// It's a float64 to handle precision. It will be 0 if no PRs have been merged.
	AvgMergeTimeSeconds float64 `json:"avg_merge_time_seconds"`
	// MergedAt is populated by the database.
	MergedAt *time.Time `json:"mergedAt,omitempty" db:"merged_at"`
}
