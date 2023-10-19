package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/mux"
	teams_rankings "github.com/jackmcguire1/riot-rankings/api/team_rankings"
	"github.com/jackmcguire1/riot-rankings/dom/riot"
	"github.com/jackmcguire1/riot-rankings/internal/sm_mux"
)

func main() {
	jsonLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(jsonLogHandler)
	riotRepo, err := riot.NewMongoRepo(
		context.Background(),
		&riot.MongoRepoParams{
			Host:                        os.Getenv("MONGO_HOST"),
			Database:                    os.Getenv("MONGO_DATABASE"),
			TeamCollection:              os.Getenv("MONGO_TEAMS_COLLECTION"),
			PlayerCollection:            os.Getenv("MONGO_PLAYERS_COLLECTION"),
			LeagueCollection:            os.Getenv("MONGO_LEAGUES_COLLECTION"),
			TournamentMappingCollection: os.Getenv("MONGO_MAPPINGS_COLLECTION"),
			GameCollection:              os.Getenv("MONGO_GAMES_COLLECTION"),
			TournamentCollection:        os.Getenv("MONGO_TOURNAMENTS_COLLECTION"),
		},
	)
	if err != nil {
		slog.
			With("error", err).
			Error("failed to init mogno repo")
		panic(err)
	}

	riotSvc := riot.NewService(&riot.Resources{Repository: riotRepo})

	teamRankingsHandler := teams_rankings.Handler{
		RiotSvc: riotSvc,
		Logger:  logger,
	}

	r := mux.NewRouter()
	r.HandleFunc("/", teamRankingsHandler.GetTeamsRankings)
	r.PathPrefix("/").HandlerFunc(teamRankingsHandler.GetTeamsRankings)

	sm := sm_mux.NewV2(r)
	lambda.Start(sm.ProxyWithContext)
}
