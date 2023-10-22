package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackmcguire1/riot-global-rankings/api/global_rankings"
	"github.com/jackmcguire1/riot-global-rankings/api/team_rankings"
	"github.com/jackmcguire1/riot-global-rankings/api/tournament_rankings"
	"github.com/jackmcguire1/riot-global-rankings/dom/riot"
	"github.com/jackmcguire1/riot-global-rankings/internal/utils"
)

func main() {
	jsonLogHandler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.Level(utils.EnvironmentWithDefaultInt("LOG_VERBOSITY", 0)),
		},
	)
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
	r := http.NewServeMux()

	globalRankingHandler := global_rankings.Handler{
		RiotSvc: riotSvc,
		Logger:  logger,
	}

	r.HandleFunc("/global_rankings", globalRankingHandler.GetGlobalRankings)

	teamRankingsHandler := team_rankings.Handler{
		RiotSvc: riotSvc,
		Logger:  logger,
	}
	r.HandleFunc("/team_rankings", teamRankingsHandler.GetTeamsRankings)

	tournamentsRankingHandler := tournament_rankings.Handler{
		RiotSvc: riotSvc,
		Logger:  logger,
	}
	r.HandleFunc("/tournament_rankings/{id:[0-9]+}", tournamentsRankingHandler.GetTournamentRankings)

	listenHost := os.Getenv("LISTEN_HOST")
	listenPort := os.Getenv("LISTEN_PORT")
	addr := fmt.Sprintf("%s:%s", listenHost, listenPort)

	logger.
		With("addr", addr).
		Info("starting http server")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGKILL)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		_ = <-signals
		err = server.Close()
		if err != nil {
			logger.
				With("error", err).
				Error("error closing server")
		}
	}()

	err = server.ListenAndServe()
	if err != nil {
		logger.
			With("error", err).
			Error("failed to listen and serve")
		panic(err)
	}
}
