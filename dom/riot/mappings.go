package riot

import "go.mongodb.org/mongo-driver/bson/primitive"

type MappingRepo interface {
	GetMappingByPlatformGameId(string) (*TournamentGameMapping, error)
	GetMappingByEsportsGameID(string) (*TournamentGameMapping, error)
	GetAllMappings() ([]*TournamentGameMapping, error)
}

type TournamentGameMapping struct {
	ID                 primitive.ObjectID `json:"_id" bson:"_id"`
	EsportsGameId      string             `json:"esportsGameId" bson:"esportsGameId"`
	PlatformGameID     string             `json:"platformGameId" bson:"platformGameId"`
	TeamMapping        map[string]string  `json:"teamMapping" bson:"teamMapping"`
	ParticipantMapping map[string]string  `json:"participantMapping" bson:"participantMapping"`
}
