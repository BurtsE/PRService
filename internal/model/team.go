package model

type TeamName string
type Team struct {
	Name    TeamName `json:"team_name"`
	Members []User   `json:"members"`
}

func (t *Team) Valid() bool {
	if len(t.Members) == 0 || t.Name == "" {
		return false
	}

	for _, member := range t.Members {
		if !member.Valid() {
			return false
		}
	}

	return true
}
