package postgres

import (
    "PRService/internal/model"
    "context"
    "github.com/jackc/pgx/v5"
)

func (r *Repository) SetUserIsActive(ctx context.Context, id model.UserID) error {
    // Mark the user as active. If the user does not exist, return an error.
    cmdTag, err := r.c.Exec(ctx, `UPDATE users SET is_active = TRUE WHERE id = $1`, id)
    if err != nil {
        return err
    }
    if cmdTag.RowsAffected() == 0 {
        // No user updated; user not found
        return pgx.ErrNoRows
    }
    return nil
}
