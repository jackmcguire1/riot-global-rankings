package riot

import "go.mongodb.org/mongo-driver/bson/primitive"

type TeamRepo interface {
	GetTeam(string) (*Team, error)
	GetAllTeams() ([]*Team, error)
	GetTeamsByIdsAndLimit([]string, int) ([]*Team, error)
	CalculateGlobalRankings(teamIDs []string, limit int) ([]*TeamRanking, error)
}

type Team struct {
	BsonID        primitive.ObjectID `bson:"_id"`
	ID            string             `json:"team_id" bson:"team_id"`
	Name          string             `json:"name" bson:"name"`
	Acronym       string             `json:"acronym" bson:"acronym"`
	Slug          string             `json:"slug" bson:"slug"`
	GlobalRanking int64              `json:"global_ranking" bson:"global_ranking"`
	TotalWins     int64              `json:"total_wins" bson:"total_wins"`
	TotalLosses   int64              `json:"total_losses" bson:"total_losses"`
	TotalTies     int64              `json:"total_ties" bson:"total_ties"`
}

// Define the struct for the result structure
type TeamRanking struct {
	TeamID    string `json:"team_id" bson:"team_id"`
	TeamCode  string `json:"team_code" bson:"team_code"`
	TeamName  string `json:"team_name" bson:"team_name"`
	Wins      int64  `json:"total_wins" bson:"total_wins"`
	Losses    int64  `json:"total_losses" bson:"total_losses"`
	Ties      int64  `json:"total_ties" bson:"total_ties"`
	Rank      int64  `json:"rank" bson:"rank"`
	EloRating int64  `json:"elo_rating" bson:"elo_rating"`
}
