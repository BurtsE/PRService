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

	// Only update name if a non-empty value is provided; otherwise preserve existing name
	userInsertQuery := `
		INSERT INTO users (id, name, team_name, is_active) 
		VALUES ($1, $2, $3, COALESCE($4, TRUE) ) 
		ON CONFLICT (id) DO 
		UPDATE 
		SET 
		  name = CASE WHEN EXCLUDED.name <> '' THEN EXCLUDED.name ELSE users.name END, 
		  is_active = COALESCE(
			EXCLUDED.is_active, users.is_active
		  ),
		  team_name = EXCLUDED.team_name
		RETURNING name
	`

	for i := range team.Members {
		err = tx.QueryRow(ctx, userInsertQuery,
			team.Members[i].ID,
			team.Members[i].Name,
			team.Name,
			team.Members[i].IsActive,
		).Scan(&team.Members[i].Name)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
