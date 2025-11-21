package model

type TeamName string
type Team struct {
	Name    TeamName `json:"team_name"`
	Members []User   `json:"members"`
}
