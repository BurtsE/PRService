package model

import "encoding/json"

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

// MarshalJSON  for correct response
func (t *Team) MarshalJSON() ([]byte, error) {
	// Define a temporary user type without TeamName
	type userDto struct {
		ID       UserID `json:"user_id"`
		Name     string `json:"username"`
		IsActive bool   `json:"is_active"`
	}

	// Convert Members to the trimmed version
	members := make([]userDto, len(t.Members))
	for i, u := range t.Members {
		members[i] = userDto{
			ID:       u.ID,
			Name:     u.Name,
			IsActive: u.IsActive,
		}
	}

	// Define an anonymous struct matching desired JSON
	tmp := struct {
		Name    TeamName  `json:"team_name"`
		Members []userDto `json:"members"`
	}{
		Name:    t.Name,
		Members: members,
	}

	return json.Marshal(tmp)
}
