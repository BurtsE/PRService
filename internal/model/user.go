package model

type UserID string
type User struct {
	ID       UserID   `json:"user_id"`
	Name     string   `json:"username"`
	TeamName TeamName `json:"team_name"`
	IsActive bool     `json:"is_active"`
}

func (u *User) Valid() bool {
	if u == nil || len(u.ID) == 0 {
		return false
	}
	return true
}
