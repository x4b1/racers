package graph

import (
	"github.com/xabi93/racers/internal/storage/postgres/ent"
)

type Race struct {
	ent.Race
}

func (r Race) Competitors() ([]*User, error) {
	competitors := make([]*User, len(r.Edges.Competitors))
	entCompetitors, err := r.Edges.CompetitorsOrErr()
	if err != nil {
		return nil, err
	}

	for i, c := range entCompetitors {
		competitors[i] = &User{*c}
	}

	return competitors, nil
}

func (Race) IsRaceResult()       {}
func (Race) IsCreateRaceResult() {}

type User struct {
	ent.User
}
