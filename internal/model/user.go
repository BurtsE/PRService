package model

type UserID string
type User struct {
	Id       UserID `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (u *User) Valid() bool {
	if u == nil || u.Id == "" {
		return false
	}
	return true
}
