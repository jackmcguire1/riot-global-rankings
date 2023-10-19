package team_rankings

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackmcguire1/riot-rankings/api"
	"github.com/jackmcguire1/riot-rankings/dom/riot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetGlobalRankings(t *testing.T) {

	mockRiotSvc := &riot.MockService{}
	mockRiotSvc.On("CalculateGlobalRankings", []string{"1", "2"}, 0).Return([]*riot.TeamRanking{
		{
			TeamID:   "1",
			TeamCode: "NIJ",
			TeamName: "NinjasInPyjamas",
			Rank:     1,
			Losses:   100,
			Wins:     1000,
		},
		{
			TeamID:   "2",
			TeamCode: "NIJ",
			TeamName: "NinjasInPyjamas",
			Rank:     1,
			Losses:   100,
			Wins:     1000,
		},
	}, nil)

	jsonLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	h := Handler{
		RiotSvc: mockRiotSvc,
		Logger:  slog.New(jsonLogHandler),
	}

	r, _ := http.NewRequest("GET", "/teams_rankings?team_ids=1,2", nil)
	w := httptest.NewRecorder()

	h.GetTeamsRankings(w, r)
	w.Body.Bytes()
	assert.EqualValues(t, http.StatusOK, w.Code)
	assert.Contains(t, string(w.Body.Bytes()), "NinjasInPyjamas")
}

func TestGetGlobalRankingsWithOutTeamIds(t *testing.T) {

	mockRiotSvc := &riot.MockService{}
	mockRiotSvc.On("CalculateGlobalRankings", []string{}, 1).Return([]*riot.TeamRanking{
		{
			TeamID:   "1",
			TeamCode: "NIJ",
			TeamName: "NinjasInPyjamas",
			Rank:     1,
			Losses:   100,
			Wins:     1000,
		},
		{
			TeamID:   "2",
			TeamCode: "NIJ",
			TeamName: "NinjasInPyjamas",
			Rank:     1,
			Losses:   100,
			Wins:     1000,
		},
	}, nil)

	jsonLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	h := Handler{
		RiotSvc: mockRiotSvc,
		Logger:  slog.New(jsonLogHandler),
	}

	r, _ := http.NewRequest("GET", "/teams_rankings", nil)
	w := httptest.NewRecorder()

	h.GetTeamsRankings(w, r)
	w.Body.Bytes()
	assert.EqualValues(t, http.StatusBadRequest, w.Code)

	var response *api.ApiError
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.EqualValues(t, response.Error, "invalid team_ids query parameter expected list of ids seperated by ','")
}

func TestGetGlobalRankingsWithInternalError(t *testing.T) {

	mockRiotSvc := &riot.MockService{}
	mockRiotSvc.On("CalculateGlobalRankings", mock.Anything, 0).Return(
		nil,
		fmt.Errorf("internal server error from mongodb"),
	)

	jsonLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	h := Handler{
		RiotSvc: mockRiotSvc,
		Logger:  slog.New(jsonLogHandler),
	}

	r, _ := http.NewRequest("GET", "/teams_rankings?team_ids=1", nil)
	w := httptest.NewRecorder()

	h.GetTeamsRankings(w, r)
	w.Body.Bytes()
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)

	var response *api.ApiError
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.EqualValues(t, response.Error, "failed to get teams rankings")
}
