package postgres

import (
	"PRService/internal/model"
	"context"
)

func (r *Repository) SetUserIsActive(ctx context.Context, user *model.User) error {
	// Mark the user as active. If the user does not exist, return an error.
	query := `
		UPDATE users SET is_active = $1 WHERE id = $2
		RETURNING id, name, team_name, is_active
	`
	err := r.c.QueryRow(ctx, query, user.IsActive, user.ID).
		Scan(&user.ID, &user.Name, &user.TeamName, &user.IsActive)
	if err != nil {
		return err
	}

	return nil
}
