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

	for i := range t.Members {
		if !t.Members[i].Valid() {
			return false
		}
	}

	return true
}
