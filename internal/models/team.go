package models

type Team struct {
	Name    string `db:"name" json:"team_name"`
	Members []User `json:"members"`
}

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}
