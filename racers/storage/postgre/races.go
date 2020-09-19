package postgre

import (
	"context"

	"github.com/xabi93/go-clean/internal/errors"
	"github.com/xabi93/go-clean/racers"
	"github.com/xabi93/go-clean/racers/service"
	"gorm.io/gorm"
)

var _ service.RacesRepository = Races{}

type Race struct {
	gorm.Model
	ID   racers.RaceID   `db:"id"`
	Name racers.RaceName `db:"name"`
	Date racers.RaceDate `db:"date"`
}

type Races struct {
	db gorm.DB
}

func (r Races) ByID(ctx context.Context, id racers.RaceID) (*racers.Race, error) {
	return nil, nil
}

func (r Races) Save(ctx context.Context, race racers.Race) error {
	err := r.db.Save(&Race{
		ID:   race.ID(),
		Name: race.Name(),
		Date: race.Date(),
	}).Error
	if err != nil {
		return errors.WrapInternalError(err, "postgre")
	}

	return nil
}
