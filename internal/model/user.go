package model

type UserID string
type User struct {
	Id       UserID `json:"user_id"`
	Name     string `json:"username"`
	IsActive bool   `json:"is_active"`
}
