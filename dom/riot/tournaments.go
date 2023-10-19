package riot

type TournamentRepo interface {
	GetTournament(string) (*Tournament, error)
	GetAllTournaments() ([]*Tournament, error)
	GetTournamentRanking(tournamentID, stageName string) ([]*TeamRanking, error)
}

type Tournament struct {
	ID        string            `json:"id" bson:"id"`
	LeagueID  string            `json:"leagueId" bson:"leagueId"`
	Name      string            `json:"name"  bson:"name"`
	Slug      string            `json:"slug"  bson:"slug"`
	Sport     string            `json:"sport"  bson:"sport"`
	StartDate string            `json:"startDate"  bson:"startDate"`
	EndDate   string            `json:"endDate"  bson:"endDate"`
	Stages    []TournamentStage `json:"stages"  bson:"stages"`
}

type TournamentStage struct {
	Name     string              `json:"name"  bson:"name"`
	Type     any                 `json:"type"  bson:"type"`
	Slug     string              `json:"slug"  bson:"slug"`
	Sections []TournamentSection `json:"sections"  bson:"sections"`
}

type TournamentSection struct {
	Name     string            `json:"name"  bson:"name"`
	Matches  []TournamentMatch `json:"matches"  bson:"matches"`
	Rankings []any             `json:"rankings"  bson:"rankings"`
}

type TournamentMatch struct {
	ID       string                  `json:"id"  bson:"id"`
	Type     string                  `json:"type"  bson:"type"`
	State    string                  `json:"state"  bson:"state"`
	Mode     string                  `json:"mode"  bson:"mode"`
	Strategy TournamentMatchStrategy `json:"strategy"  bson:"strategy"`
	Teams    []TournamentMatchTeam   `json:"teams" bson:"teams"`
	Games    []TournamentMatchGame   `json:"games" bson:"games"`
}

type TournamentMatchGame struct {
	ID     string                    `json:"id" bson:"id"`
	State  string                    `json:"state" bson:"state"`
	Number int                       `json:"number" bson:"number"`
	Teams  []TournamentMatchGameTeam `json:"teams" bson:"teams"`
}

type TournamentMatchGameTeam struct {
	ID     string `json:"id" bson:"id"`
	Side   string `json:"side" bson:"side"`
	Result struct {
		Outcome string `json:"outcome" bson:"outcome"`
	} `json:"result" bson:"result"`
}

type TournamentMatchStrategy struct {
	Type  string `json:"type" bson:"type"`
	Count int    `json:"count" bson:"count"`
}

type TournamentMatchTeam struct {
	ID      string                      `json:"id" bson:"id"`
	Side    string                      `json:"side" bson:"side"`
	Record  TournamentMatchTeamRecord   `json:"record" bson:"record"`
	Result  TournamnetMatchTeamResults  `json:"result" bson:"result"`
	Players []TournamentMatchTeamPlayer `json:"players" bson:"players"`
}

type TournamentMatchTeamPlayer struct {
	ID   string `json:"id" bson:"id"`
	Role string `json:"role" bson:"role"`
}

type TournamnetMatchTeamResults struct {
	Outcome  string `json:"outcome" bson:"outcome"`
	GameWins int    `json:"gameWins" bson:"gameWins"`
}

type TournamentMatchTeamRecord struct {
	Wins   int `json:"wins" bson:"wins"`
	Losses int `json:"losses" bson:"losses"`
	Ties   int `json:"ties" bson:"ties"`
}
