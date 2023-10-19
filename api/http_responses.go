package api

type ApiError struct {
	Error string `json:"error"`
}

type TeamResponse struct {
	TeamID      string `json:"team_id"`
	TeamCode    string `json:"team_code"`
	TeamName    string `json:"team_name"`
	Rank        int64  `json:"rank"`
	TotalWins   int64  `json:"total_wins,omitempty"`
	TotalLosses int64  `json:"total_losses,omitempty"`
	TotalTies   int64  `json:"total_ties,omitempty"`
}
