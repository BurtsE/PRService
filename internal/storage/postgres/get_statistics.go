package postgres

import (
	"PRService/internal/model"
	"context"
)

func (r *Repository) GetStatistic(ctx context.Context) (model.Statistic, error) {
	query := `
		WITH pr_stats AS (
			SELECT
				COUNT(*) AS total_prs,
				COUNT(*) FILTER (WHERE status = 'OPEN') AS open_prs,
				COUNT(*) FILTER (WHERE status = 'MERGED') AS merged_prs,
				AVG(EXTRACT(EPOCH FROM (merged_at - created_at))) AS avg_merge_seconds
			FROM pull_requests
		),
		user_stats AS (
			SELECT
				COUNT(*) AS total_users,
				COUNT(*) FILTER (WHERE is_active = TRUE) AS active_users
			FROM users
		),
		team_stats AS (
			SELECT COUNT(*) AS total_teams FROM teams
		)
		SELECT
			pr_stats.total_prs,
			pr_stats.open_prs,
			pr_stats.merged_prs,
			COALESCE(pr_stats.avg_merge_seconds, 0) as avg_merge_seconds,
			user_stats.total_users,
			user_stats.active_users,
			team_stats.total_teams
		FROM pr_stats, user_stats, team_stats;
	`

	var stats model.Statistic
	err := r.c.QueryRow(ctx, query).Scan(&stats.TotalPRs, &stats.OpenPRs, &stats.MergedPRs, &stats.AvgMergeTimeSeconds, &stats.TotalUsers, &stats.ActiveUsers, &stats.TotalTeams)
	return stats, err
}
