package riot

type PlayerRepo interface {
	GetPlayer(string) (*Player, error)
	GetAllPlayers() ([]*Player, error)
}

type Player struct {
	PlayerID   string `json:"player_id" bson:"player_id"`
	Handle     string `json:"handle" bson:"handle"`
	FirstName  string `json:"first_name" bson:"first_name"`
	LastName   string `json:"last_name" bson:"last_name"`
	HomeTeamID string `json:"home_team_id" bson:"home_team_id"`
}
