package team_rankings

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/jackmcguire1/riot-global-rankings/api"
	"github.com/jackmcguire1/riot-global-rankings/dom/riot"
	"github.com/jackmcguire1/riot-global-rankings/internal/utils"
)

type Handler struct {
	RiotSvc riot.Service
	Logger  *slog.Logger
}

func (h *Handler) GetTeamsRankings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(utils.ToJsonBytes(api.ApiError{Error: "HTTP METHOD not supported"}))
		return
	}

	teamIds := []string{}
	if teamIdsQueryParam := r.URL.Query().Get("team_ids"); teamIdsQueryParam != "" {
		teamIds = strings.Split(teamIdsQueryParam, ",")
		h.Logger.
			With("limit-query-parameter", teamIdsQueryParam).
			With("new-limit", teamIds).
			Debug("client has specified limit query parameter")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(utils.ToJsonBytes(api.ApiError{Error: "invalid team_ids query parameter expected list of ids seperated by ','"}))
		return
	}

	teams, err := h.RiotSvc.CalculateGlobalRankings(teamIds, 0)
	if err != nil {
		h.Logger.
			With("error", err).
			Error("failed to get global rankings")

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(utils.ToJsonBytes(api.ApiError{Error: "failed to get teams rankings"}))
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
