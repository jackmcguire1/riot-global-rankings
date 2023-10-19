package riot

type Repository interface {
	TeamRepo
	MappingRepo
	LeagueRepo
	PlayerRepo
	TournamentRepo
}
