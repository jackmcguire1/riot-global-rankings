package riot

type LeagueRepo interface {
	GetLeague(string) (*League, error)
	GetAllLeagues() ([]*League, error)
}

type League struct {
	ID              string                `json:"id" bson:"id"`
	Name            string                `json:"name" bson:"name"`
	Slug            string                `json:"slug" bson:"slug"`
	Sport           string                `json:"sport" bson:"sport"`
	Image           string                `json:"image" bson:"image"`
	LightImage      string                `json:"lightImage" bson:"lightImage"`
	DarkImage       string                `json:"darkImage" bson:"darkImage"`
	Region          string                `json:"region" bson:"region"`
	Priority        int64                 `json:"priority" bson:"priority"`
	DisplayPriority LeagueDisplayPriority `json:"displayPriority" bson:"displayPriority"`
	Tournaments     []LeagueTournament    `json:"tournaments" bson:"tournaments"`
}

type LeagueTournament struct {
	ID string `json:"id" bson:"id"`
}

type LeagueDisplayPriority struct {
	Position int    `json:"position" bson:"position"`
	Status   string `json:"status" bson:"status"`
}
