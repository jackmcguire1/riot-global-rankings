package global_rankings

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackmcguire1/riot-global-rankings/api"
	"github.com/jackmcguire1/riot-global-rankings/dom/riot"
	"github.com/jackmcguire1/riot-global-rankings/internal/utils"
)

type Handler struct {
	RiotSvc riot.Service
	Logger  *slog.Logger
}

func (h *Handler) GetGlobalRankings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(utils.ToJsonBytes(api.ApiError{Error: "HTTP METHOD not supported"}))
		return
	}

	limit := 20
	if limitQueryParam := r.URL.Query().Get("number_of_teams"); limitQueryParam != "" {
		limit, _ = strconv.Atoi(limitQueryParam)
		h.Logger.
			With("limit-query-parameter", limitQueryParam).
			With("new-limit", limit).
			Debug("client has specified limit query parameter")
	}

	teams, err := h.RiotSvc.CalculateGlobalRankings([]string{}, limit)
	if err != nil {
		h.Logger.
			With("error", err).
			Error("failed to get global rankings")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ToJsonBytes(api.ApiError{Error: "failed to get global rankings"}))
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
			TotalTies:   team.Ties,
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Write(utils.PrettyPrintJson(apiResponse))
	return
}
