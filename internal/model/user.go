package model

type UserID string
type User struct {
	Id       UserID `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (u *User) Valid() bool {
	if u == nil || len(u.Id) == 0 {
		return false
	}
	return true
}
