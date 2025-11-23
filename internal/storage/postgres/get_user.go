package postgres

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
	"database/sql"
	"errors"
)

func (r *Repository) GetUser(ctx context.Context, id model.UserID) (model.User, error) {
	var user model.User
	query := `SELECT id, name, team_name, is_active FROM users WHERE id = $1`

	err := r.c.QueryRow(ctx, query, id).Scan(&user.ID, &user.Name, &user.TeamName, &user.IsActive)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, service.ErrResourceNotFound
	}
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
