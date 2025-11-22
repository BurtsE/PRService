package postgres

import (
	"PRService/internal/model"
	"context"
)

func (r *Repository) GetTeam(ctx context.Context, teamName model.TeamName) (*model.Team, error) {
	// Fetch members for the team
	rows, err := r.c.Query(ctx, `
		SELECT u.id, u.name, u.is_active
		FROM team_members tm
		JOIN users u ON u.id = tm.user_id
		WHERE tm.team_name = $1
	`, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]model.User, 0)
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.Id, &u.Name, &u.IsActive); err != nil {
			return nil, err
		}
		members = append(members, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	team := &model.Team{
		Name:    teamName,
		Members: members,
	}
	return team, nil
}
