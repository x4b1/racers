package service

type ActionName string

func (an ActionName) String() string {
	return string(an)
}

func NewService(races RacesService, teams TeamsService) *Service {
	s := Service{
		races,
		teams,
	}

	return &s
}

type Service struct {
	Races RacesService
	Teams TeamsService
}
