package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackmcguire1/riot-rankings/dom/riot"
	"github.com/jackmcguire1/riot-rankings/internal/utils"
)

func main() {
	hndler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(hndler)

	mongoParams := &riot.MongoRepoParams{
		Host:                        "",
		Database:                    "riot",
		TeamCollection:              "teams",
		PlayerCollection:            "players",
		LeagueCollection:            "leagues",
		TournamentMappingCollection: "mappings",
		TournamentCollection:        "tournaments",
	}
	repo, err := riot.NewMongoRepo(context.Background(), mongoParams)
	if err != nil {
		logger.With("error", err).Error("failed to init mongo repo")
		panic(err)
	}

	svc := riot.NewService(&riot.Resources{Repository: repo})

	results, err := svc.CalculateGlobalRankings([]string{"105515219038427019", "103935530333072898", "98767991954244555"}, 2)
	if err != nil {
		logger.With("error", err).Error("failed to calculate global rankings")
		panic(err)
	}

	logger.
		With("len-teams", len(results)).
		Info("got global rankings")

	globalTeamRankings, err := svc.CalculateGlobalRankings([]string{}, 20)
	if err != nil {
		logger.With("error", err).Error("failed to get all teams")
		panic(err)
	}
	logger.
		With("len-teams", len(globalTeamRankings)).
		With("team-info", utils.ToJSON(globalTeamRankings)).
		Info("got global rankings")
}
