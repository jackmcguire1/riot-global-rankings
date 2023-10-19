package riot

import "github.com/stretchr/testify/mock"

type MockService struct {
	mock.Mock
}

func (svc *MockService) GetAllTeams() (teams []*Team, err error) {
	args := svc.Called()
	if args.Get(0) != nil {
		teams = args.Get(0).([]*Team)
	}
	return teams, args.Error(1)
}

func (svc *MockService) CalculateGlobalRankings(teamIDs []string, limit int) (t []*TeamRanking, err error) {
	args := svc.Called(teamIDs, limit)
	if args.Get(0) != nil {
		t = args.Get(0).([]*TeamRanking)
	}
	return t, args.Error(1)
}

func (svc *MockService) GetMappingByPlatformGameId(s string) (mapping *TournamentGameMapping, err error) {
	args := svc.Called()
	if args.Get(0) != nil {
		mapping = args.Get(0).(*TournamentGameMapping)
	}
	return mapping, args.Error(1)
}

func (svc *MockService) GetMappingByEsportsGameID(s string) (mapping *TournamentGameMapping, err error) {
	args := svc.Called()
	if args.Get(0) != nil {
		mapping = args.Get(0).(*TournamentGameMapping)
	}
	return mapping, args.Error(1)
}

func (svc *MockService) GetAllMappings() (mappings []*TournamentGameMapping, err error) {
	args := svc.Called()
	if args.Get(0) != nil {
		mappings = args.Get(0).([]*TournamentGameMapping)
	}
	return mappings, args.Error(1)
}

func (svc *MockService) GetLeague(s string) (league *League, err error) {
	args := svc.Called()
	if args.Get(0) != nil {
		league = args.Get(0).(*League)
	}
	return league, args.Error(1)
}

func (svc *MockService) GetAllLeagues() (leagues []*League, err error) {
	args := svc.Called()
	if args.Get(0) != nil {
		leagues = args.Get(0).([]*League)
	}
	return leagues, args.Error(1)
}

func (svc *MockService) GetPlayer(s string) (player *Player, err error) {
	args := svc.Called()
	if args.Get(0) != nil {
		player = args.Get(0).(*Player)
	}
	return player, args.Error(1)
}

func (svc *MockService) GetAllPlayers() (players []*Player, err error) {
	args := svc.Called()
	if args.Get(0) != nil {
		players = args.Get(0).([]*Player)
	}
	return players, args.Error(1)
}

func (svc *MockService) GetTournament(s string) (tourny *Tournament, err error) {
	args := svc.Called()
	if args.Get(0) != nil {
		tourny = args.Get(0).(*Tournament)
	}
	return tourny, args.Error(1)
}

func (svc *MockService) GetAllTournaments() (tournys []*Tournament, err error) {
	args := svc.Called()
	if args.Get(0) != nil {
		tournys = args.Get(0).([]*Tournament)
	}
	return tournys, args.Error(1)
}

func (svc *MockService) GetTeam(id string) (team *Team, err error) {
	args := svc.Called(id)
	if args.Get(0) != nil {
		team = args.Get(0).(*Team)
	}
	return team, args.Error(1)
}

func (svc *MockService) GetTeamsByIdsAndLimit(teamIDs []string, limit int) (teams []*Team, err error) {
	args := svc.Called(teamIDs, limit)
	if args.Get(0) != nil {
		teams = args.Get(0).([]*Team)
	}
	return teams, args.Error(1)
}

func (svc *MockService) GetTournamentRanking(tournamentID string, stage string) (teams []*TeamRanking, err error) {
	args := svc.Called(tournamentID, stage)
	if args.Get(0) != nil {
		teams = args.Get(0).([]*TeamRanking)
	}
	return teams, args.Error(1)
}
