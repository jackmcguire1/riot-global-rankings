package global_rankings

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
	mockRiotSvc.On("CalculateGlobalRankings", mock.Anything, 20).Return([]*riot.TeamRanking{
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
			TeamCode: "GNC",
			TeamName: "Giga Ninjas",
			Rank:     2,
			Losses:   10,
			Wins:     10,
		},
		{
			TeamID:   "3",
			TeamCode: "PIP",
			TeamName: "Plumps in Pjs",
			Rank:     3,
			Losses:   10,
			Wins:     0,
		},
	}, nil)

	jsonLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	h := Handler{
		RiotSvc: mockRiotSvc,
		Logger:  slog.New(jsonLogHandler),
	}

	r, _ := http.NewRequest("GET", "/global_rankings", nil)
	w := httptest.NewRecorder()

	h.GetGlobalRankings(w, r)
	w.Body.Bytes()
	assert.EqualValues(t, http.StatusOK, w.Code)
	assert.Contains(t, string(w.Body.Bytes()), "NinjasInPyjamas")
}

func TestGetGlobalRankingsWithLimit(t *testing.T) {

	mockRiotSvc := &riot.MockService{}
	mockRiotSvc.On("CalculateGlobalRankings", mock.Anything, 1).Return([]*riot.TeamRanking{
		{
			TeamID:   "1",
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

	r, _ := http.NewRequest("GET", "/global_rankings?number_of_teams=1", nil)
	w := httptest.NewRecorder()

	h.GetGlobalRankings(w, r)
	w.Body.Bytes()
	assert.EqualValues(t, http.StatusOK, w.Code)

	var response []*api.TeamResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response, 1)
	assert.NotContains(t, string(w.Body.Bytes()), "Cloud0")
}

func TestGetGlobalRankingsWithInternalError(t *testing.T) {

	mockRiotSvc := &riot.MockService{}
	mockRiotSvc.On("CalculateGlobalRankings", mock.Anything, 1).Return(
		nil,
		fmt.Errorf("internal server error from mongodb"),
	)

	jsonLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	h := Handler{
		RiotSvc: mockRiotSvc,
		Logger:  slog.New(jsonLogHandler),
	}

	r, _ := http.NewRequest("GET", "/global_rankings?number_of_teams=1", nil)
	w := httptest.NewRecorder()

	h.GetGlobalRankings(w, r)
	w.Body.Bytes()
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)

	var response *api.ApiError
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.EqualValues(t, response.Error, "failed to get global rankings")
}
