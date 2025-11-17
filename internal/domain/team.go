package domain

type Team struct {
	Name    string  `json:"team_name" db:"name"`
	Members []*User `json:"members"`
}

type TeamStats struct {
	Name  string              `json:"team_name"`
	Users []*UserStats        `json:"users"`
	PRs   []*PullRequestStats `json:"open_prs"`
}
