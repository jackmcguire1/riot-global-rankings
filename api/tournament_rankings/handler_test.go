package tournament_rankings

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

func TestGetTournamentRankings(t *testing.T) {

	mockRiotSvc := &riot.MockService{}
	mockRiotSvc.On("GetTournamentRanking", mock.Anything, mock.Anything).Return([]*riot.TeamRanking{
		{
			TeamID:   "1",
			TeamCode: "NIJ",
			TeamName: "NinjasInPyjamas",
			Rank:     1,
			Losses:   0,
			Wins:     10,
		},
		{
			TeamID:   "1",
			TeamCode: "NIJ",
			TeamName: "NinjasInPyjamas",
			Rank:     2,
			Losses:   5,
			Wins:     5,
		},
		{
			TeamID:   "1",
			TeamCode: "NIJ",
			TeamName: "NinjasInPyjamas",
			Rank:     3,
			Losses:   0,
			Wins:     10,
		},
	}, nil)

	jsonLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	h := Handler{
		RiotSvc: mockRiotSvc,
		Logger:  slog.New(jsonLogHandler),
	}

	r, _ := http.NewRequest("GET", "/tournament_rankings/1", nil)
	w := httptest.NewRecorder()

	h.GetTournamentRankings(w, r)
	w.Body.Bytes()
	assert.EqualValues(t, http.StatusOK, w.Code)
	assert.Contains(t, string(w.Body.Bytes()), "NinjasInPyjamas")
}

func TestGetGlobalRankingsWithLimit(t *testing.T) {

	mockRiotSvc := &riot.MockService{}
	mockRiotSvc.On("GetTournamentRanking", mock.Anything, "test").Return([]*riot.TeamRanking{
		{
			TeamID:   "1",
			TeamCode: "NIJ",
			TeamName: "NinjasInPyjamas",
			Rank:     1,
			Losses:   0,
			Wins:     10,
		},
	}, nil)

	jsonLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	h := Handler{
		RiotSvc: mockRiotSvc,
		Logger:  slog.New(jsonLogHandler),
	}

	r, _ := http.NewRequest("GET", "/tournament_rankings/1?stage=test", nil)
	w := httptest.NewRecorder()

	h.GetTournamentRankings(w, r)
	w.Body.Bytes()
	assert.EqualValues(t, http.StatusOK, w.Code)

	var response []*api.TeamResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(t, response, 1)
	assert.NotContains(t, string(w.Body.Bytes()), "Cloud0")
}

func TestGetGlobalRankingsWithInternalError(t *testing.T) {

	mockRiotSvc := &riot.MockService{}
	mockRiotSvc.On("GetTournamentRanking", mock.Anything, mock.Anything).Return(
		nil,
		fmt.Errorf("internal server error from mongodb"),
	)

	jsonLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	h := Handler{
		RiotSvc: mockRiotSvc,
		Logger:  slog.New(jsonLogHandler),
	}

	r, _ := http.NewRequest("GET", "/tournament_rankings/1?stage=0", nil)
	w := httptest.NewRecorder()

	h.GetTournamentRankings(w, r)
	w.Body.Bytes()
	assert.EqualValues(t, http.StatusInternalServerError, w.Code)

	var response *api.ApiError
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.EqualValues(t, response.Error, "failed to get tournament rankings")
}
