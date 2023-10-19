package tournament_rankings

import (
	"log/slog"
	"net/http"
	"path"

	"github.com/jackmcguire1/riot-rankings/api"
	"github.com/jackmcguire1/riot-rankings/dom/riot"
	"github.com/jackmcguire1/riot-rankings/internal/utils"
)

type Handler struct {
	RiotSvc riot.Service
	Logger  *slog.Logger
}

func (h *Handler) GetTournamentRankings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(utils.ToJsonBytes(api.ApiError{Error: "HTTP METHOD not supported"}))
		return
	}

	limit := 20
	stage := r.URL.Query().Get("stage")
	if stage != "" {
		h.Logger.
			With("limit-query-parameter", stage).
			With("new-limit", limit).
			Info("client has specified tournament stage")
	}

	tournamentID := path.Base(r.URL.Path)
	if tournamentID == "" {
		h.Logger.
			With("error", "tournament-id not specified").
			Error("tournament-id not specified")

		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.ToJsonBytes(api.ApiError{Error: "invalid tournament id path parameter"}))
		return
	}

	teams, err := h.RiotSvc.GetTournamentRanking(tournamentID, stage)
	if err != nil {
		h.Logger.
			With("tournament-id", tournamentID).
			With("error", err).
			Error("failed to get tournament rankings")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ToJsonBytes(api.ApiError{Error: "failed to get tournament rankings"}))
		return
	}

	var apiResponse []*api.TeamResponse
	for _, team := range teams {
		apiResponse = append(apiResponse, &api.TeamResponse{
			TeamID:      team.TeamID,
			TeamName:    team.TeamName,
			TeamCode:    team.TeamCode,
			Rank:        team.Rank,
			TotalWins:   team.Wins,
			TotalLosses: team.Losses,
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Write(utils.PrettyPrintJson(apiResponse))
	return
}
