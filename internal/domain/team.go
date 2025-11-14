package domain

type Team struct {
	Name    string  `json:"team_name" db:"name"`
	Members []*User `json:"members"`
}
