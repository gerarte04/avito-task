package domain

type User struct {
	ID       string `json:"user_id" db:"id"`
	Name     string `json:"username" db:"name"`
	IsActive bool   `json:"is_active" db:"is_active"`

	TeamName string `json:"team_name,omitempty"`
}
