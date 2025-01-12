package model

type App struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	OrgID     string `json:"org_id"`
	APIToken  string `json:"api_token"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	IsRemoved bool   `json:"is_removed"`
}
