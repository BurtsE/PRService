package postgres

import (
	"PRService/internal/model"
	"context"
)

func (r *Repository) CreateTeam(ctx context.Context, team *model.Team) error {

	tx, err := r.c.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err = tx.Exec(ctx, `INSERT INTO teams (name) VALUES ($1)`, team.Name); err != nil {
		return err
	}

	// Ensure members exist in users table and add to team_members
	for _, u := range team.Members {
		if _, err = tx.Exec(ctx, `INSERT INTO users (id, name, is_active) VALUES ($1, $2, COALESCE($3, TRUE))
                                   ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name`, u.Id, u.Name, u.IsActive); err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, `INSERT INTO team_members (team_name, user_id) VALUES ($1, $2)
                                   ON CONFLICT (team_name, user_id) DO NOTHING`, team.Name, u.Id); err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
