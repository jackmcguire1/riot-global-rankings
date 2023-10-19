package riot

type Service interface {
	TeamRepo
	MappingRepo
	LeagueRepo
	PlayerRepo
	TournamentRepo
}

type Resources struct {
	Repository
}

type RiotService struct {
	*Resources
}

func NewService(r *Resources) *RiotService {
	return &RiotService{Resources: r}
}
