package riot

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/jackmcguire1/riot-rankings/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	TeamRepo
	MappingRepo
	LeagueRepo
	PlayerRepo
	TournamentRepo

	TeamCollection              *mongo.Collection
	PlayerCollection            *mongo.Collection
	LeagueCollection            *mongo.Collection
	TournamentMappingCollection *mongo.Collection
	TournamentCollection        *mongo.Collection
	GameCollection              *mongo.Collection
}

type MongoRepoParams struct {
	Host                        string
	Database                    string
	TeamCollection              string
	PlayerCollection            string
	LeagueCollection            string
	TournamentMappingCollection string
	GameCollection              string
	TournamentCollection        string
}

func NewMongoRepo(ctx context.Context, params *MongoRepoParams) (*MongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(params.Host))
	if err != nil {
		return nil, err
	}
	database := client.Database(params.Database)

	return &MongoRepository{
		TeamCollection:              database.Collection(params.TeamCollection),
		PlayerCollection:            database.Collection(params.PlayerCollection),
		LeagueCollection:            database.Collection(params.LeagueCollection),
		TournamentMappingCollection: database.Collection(params.TournamentMappingCollection),
		TournamentCollection:        database.Collection(params.TournamentCollection),
		GameCollection:              database.Collection(params.GameCollection),
	}, nil
}

func (repo *MongoRepository) GetTeam(ID string) (t *Team, err error) {

	filter := bson.M{"team_id": ID}
	t = &Team{}
	err = repo.getOne(filter, repo.TeamCollection, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (repo *MongoRepository) GetPlayer(ID string) (p *Player, err error) {
	filter := bson.M{"player_id": ID}
	p = &Player{}
	err = repo.getOne(filter, repo.PlayerCollection, &p)
	if err != nil {
		return nil, err
	}
	return
}

func (repo *MongoRepository) GetLeague(ID string) (l *League, err error) {
	filter := bson.M{"id": ID}
	l = &League{}
	err = repo.getOne(filter, repo.LeagueCollection, &l)
	if err != nil {
		return nil, err
	}
	return
}

func (repo *MongoRepository) GetMappingByPlatformGameId(ID string) (t *TournamentGameMapping, err error) {
	filter := bson.M{"platformGameId": ID}
	t = &TournamentGameMapping{}
	err = repo.getOne(filter, repo.TournamentMappingCollection, &t)
	if err != nil {
		return nil, err
	}
	return
}

func (repo *MongoRepository) GetMappingByEsportsGameID(ID string) (p *TournamentGameMapping, err error) {
	filter := bson.M{"esportsGameId": ID}
	p = &TournamentGameMapping{}
	err = repo.getOne(filter, repo.TournamentMappingCollection, &p)
	if err != nil {
		return nil, err
	}
	return
}

func (repo *MongoRepository) GetAllTeams() ([]*Team, error) {
	filter := bson.M{}
	slice := []*Team{}

	err := repo.searchCollection(filter, repo.TeamCollection, &slice)
	if err != nil {
		return nil, err
	}

	return slice, nil
}

func (repo *MongoRepository) GetAllMappings() ([]*TournamentGameMapping, error) {
	filter := bson.M{}
	slice := []*TournamentGameMapping{}

	err := repo.searchCollection(filter, repo.TournamentMappingCollection, &slice)
	if err != nil {
		return nil, err
	}

	return slice, nil
}

func (repo *MongoRepository) GetAllPlayers() ([]*Player, error) {
	filter := bson.M{}
	slice := []*Player{}

	err := repo.searchCollection(filter, repo.PlayerCollection, &slice)
	if err != nil {
		return nil, err
	}

	return slice, nil
}

func (repo *MongoRepository) GetAllTournaments() ([]*Tournament, error) {
	filter := bson.M{}
	slice := []*Tournament{}

	err := repo.searchCollection(filter, repo.TournamentCollection, &slice)
	if err != nil {
		return nil, err
	}

	return slice, nil
}

func (repo *MongoRepository) GetAllLeagues() ([]*League, error) {
	filter := bson.M{}
	slice := []*League{}

	err := repo.searchCollection(filter, repo.LeagueCollection, &slice)
	if err != nil {
		return nil, err
	}

	return slice, nil
}

func (repo *MongoRepository) GetTeamsByIdsAndLimit(teamIDs []string, limit int) ([]*Team, error) {

	var teamIdsMatchStage bson.M = bson.M{
		"$match": bson.M{},
	}

	var limitQ = bson.M{}
	if limit > 0 {
		limitQ = bson.M{
			"$limit": limit,
		}
	}

	if len(teamIDs) > 0 {
		teamIdsMatchStage = bson.M{
			"$match": bson.M{
				"team_id": bson.M{
					"$in": teamIDs,
				},
			},
		}
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"global_ranking": bson.M{"$exists": true},
			},
		},
		teamIdsMatchStage,
		{
			"$sort": bson.M{
				"global_ranking": 1,
			},
		},
		limitQ,
	}

	cursor, err := repo.TeamCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	teams := []*Team{}
	err = cursor.All(context.Background(), &teams)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func (repo *MongoRepository) CalculateGlobalRankings(teamIDs []string, limit int) ([]*TeamRanking, error) {
	// Your aggregation pipeline to calculate team rankings using Elo rating system
	pipeline := []bson.M{
		{
			"$unwind": "$stages",
		},
		{
			"$unwind": "$stages.sections",
		},
		{
			"$unwind": "$stages.sections.matches",
		},
		{
			"$unwind": "$stages.sections.matches.games",
		},
		{
			"$match": bson.M{
				"stages.sections.matches.games.state":   "completed",
				"stages.sections.matches.mode":          "classic",
				"stages.sections.matches.strategy.type": "bestOf",
			},
		},
		{
			"$unwind": "$stages.sections.matches.games.teams",
		},
		{
			"$group": bson.M{
				"_id":     "$stages.sections.matches.games.teams.id",
				"team_id": bson.M{"$first": "$stages.sections.matches.games.teams.id"}, // Keep the team_id as an accumulator object
				"total_wins": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$stages.sections.matches.games.teams.result.outcome", "win"}},
							1,
							0,
						},
					},
				},
				"total_losses": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$stages.sections.matches.games.teams.result.outcome", "loss"}},
							1,
							0,
						},
					},
				},
				"total_ties": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$stages.sections.matches.games.teams.result.outcome", "tie"}},
							1,
							0,
						},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"elo_rating": bson.M{
					"$subtract": []interface{}{
						"$total_wins",
						"$total_losses",
					},
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "teams",     // Replace with the actual name of the team collection
				"localField":   "team_id",   // Field in the current collection
				"foreignField": "team_id",   // Field in the team collection to match
				"as":           "team_data", // Alias for the joined data
			},
		},
		{
			"$unwind": "$team_data",
		},
		{
			"$sort": bson.M{
				"elo_rating": -1,
			},
		},
		{
			"$project": bson.M{
				"team_id":      "$team_id",
				"total_wins":   "$total_wins",
				"total_losses": "$total_losses",
				"total_ties":   "$total_ties",
				"elo_rating":   "$elo_rating",
				"team_name":    "$team_data.name",
				"team_code":    "$team_data.acronym",
			},
		},
	}

	// Execute the aggregation
	cursor, err := repo.TournamentCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Define a slice to store the results

	var teamRankings = []*TeamRanking{}
	// Iterate through the results and decode them into the TeamRanking struct

	err = cursor.All(context.Background(), &teamRankings)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var tmpRankings = []*TeamRanking{}
		err = cursor.All(context.Background(), &tmpRankings)
		if err != nil {
			return nil, err
		}

		teamRankings = append(teamRankings, tmpRankings...)
	}

	teamIdentifiers := map[string]int{}
	for rank, ranking := range teamRankings {
		rank++
		ranking.Rank = int64(rank)
		if slices.Contains(teamIDs, ranking.TeamID) {
			teamIdentifiers[ranking.TeamID] = rank
		}
	}

	if len(teamIDs) > 0 {
		var results []*TeamRanking
		for _, teamID := range teamIDs {
			rank, exists := teamIdentifiers[teamID]
			if exists {
				results = append(results, teamRankings[rank-1]) // Adjust rank to slice index
			}
		}
		return results, nil
	}

	if limit > 0 {

		index := 0
		if limit > len(teamRankings)-1 {
			index = len(teamRankings) - 1
		} else {
			index = limit - 1
		}
		results := make([]*TeamRanking, index)
		copy(results, teamRankings[:index])
		return results, nil
	}

	return teamRankings, nil
}

func (repo *MongoRepository) GetTournamentRanking(tournamentID, stageName string) ([]*TeamRanking, error) {

	tournament, err := repo.GetTournament(tournamentID)
	if err != nil {
		return nil, err
	}

	query := []bson.M{
		{
			"$match": bson.M{"startDate": bson.M{"$lt": tournament.StartDate}},
		},
		{
			"$unwind": "$stages",
		},
	}
	tournamentStages := map[string]*TournamentStage{}
	if stageName != "" {
		for _, stage := range tournament.Stages {
			tournamentStages[stage.Name] = &stage
			break
		}
	}

	_, ok := tournamentStages[stageName]
	if stageName != "" && !ok {
		return nil, fmt.Errorf("stage does not exist within tournament")
	}

	if ok {
		query = []bson.M{
			{
				"$match": bson.M{"$id": tournamentID},
			},
			{
				"$unwind": "$stages",
			},
			{
				"stages.name": bson.M{"$in": utils.MapKeys(tournamentStages)},
			},
		}
	}

	// Your aggregation pipeline to calculate team rankings using Elo rating system
	pipeline := append([]bson.M{}, query...)
	pipeline = append(pipeline,
		bson.M{
			"$unwind": "$stages.sections",
		},
		bson.M{
			"$unwind": "$stages.sections.matches",
		},
		bson.M{
			"$unwind": "$stages.sections.matches.games",
		},
		bson.M{
			"$match": bson.M{
				"stages.sections.matches.games.state":   "completed",
				"stages.sections.matches.mode":          "classic",
				"stages.sections.matches.strategy.type": "bestOf",
			},
		},
		bson.M{
			"$unwind": "$stages.sections.matches.games.teams",
		},
		bson.M{
			"$group": bson.M{
				"_id":     "$stages.sections.matches.games.teams.id",
				"team_id": bson.M{"$first": "$stages.sections.matches.games.teams.id"}, // Keep the team_id as an accumulator object
				"total_wins": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$stages.sections.matches.games.teams.result.outcome", "win"}},
							1,
							0,
						},
					},
				},
				"total_ties": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$stages.sections.matches.games.teams.result.outcome", "tie"}},
							1,
							0,
						},
					},
				},
				"total_losses": bson.M{
					"$sum": bson.M{
						"$cond": []interface{}{
							bson.M{"$eq": []interface{}{"$stages.sections.matches.games.teams.result.outcome", "loss"}},
							1,
							0,
						},
					},
				},
			},
		},
		bson.M{
			"$addFields": bson.M{
				"elo_rating": bson.M{
					"$subtract": []interface{}{
						"$total_wins",
						"$total_losses",
					},
				},
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         "teams",     // Replace with the actual name of the team collection
				"localField":   "team_id",   // Field in the current collection
				"foreignField": "team_id",   // Field in the team collection to match
				"as":           "team_data", // Alias for the joined data
			},
		},
		bson.M{
			"$unwind": "$team_data",
		},
		bson.M{
			"$sort": bson.M{
				"elo_rating": -1,
			},
		},
		bson.M{
			"$project": bson.M{
				"team_id":      "$team_id",
				"total_wins":   "$total_wins",
				"total_losses": "$total_losses",
				"total_ties":   "$total_ties",
				"elo_rating":   "$elo_rating",
				"team_name":    "$team_data.name",
				"team_code":    "$team_data.acronym",
			},
		},
	)

	// Execute the aggregation
	cursor, err := repo.TournamentCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Define a slice to store the results

	var teamRankings = []*TeamRanking{}
	// Iterate through the results and decode them into the TeamRanking struct

	for cursor.Next(context.Background()) {

		var tmpTeamRankings = []*TeamRanking{}
		err = cursor.All(context.Background(), &tmpTeamRankings)
		if err != nil {
			return nil, err
		}
		teamRankings = append(teamRankings, tmpTeamRankings...)
	}

	for i, ranking := range teamRankings {
		i++
		ranking.Rank = int64(i)
	}

	return teamRankings, nil
}

func (repo *MongoRepository) GetTournament(ID string) (t *Tournament, err error) {
	filter := bson.M{"id": ID}
	t = &Tournament{}
	err = repo.getOne(filter, repo.TournamentCollection, &t)
	if err != nil {
		return nil, err
	}
	return
}

func (repo *MongoRepository) searchCollection(filter bson.M, collection *mongo.Collection, i interface{}) error {
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return err
	}
	if cursor.Err() != nil {
		return err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), i)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MongoRepository) getOne(filter interface{}, collection *mongo.Collection, i interface{}) error {
	res := collection.FindOne(context.Background(), filter, nil)
	if res.Err() != nil {
		if strings.Contains(res.Err().Error(), "no documents in result") {
			return utils.ErrNotFound
		}
		return res.Err()
	}

	err := res.Decode(i)
	if err != nil {
		err = fmt.Errorf("failed to umarshal bson user document err:%w", err)
		return err
	}

	return nil
}
