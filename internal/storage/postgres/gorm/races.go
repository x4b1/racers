package gorm

import (
	"context"
	"time"

	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/errors"
	"github.com/xabi93/racers/internal/service"
	"gorm.io/gorm"
)

var _ service.RacesRepository = Races{}

type raceDB struct {
	ID   racers.RaceID   `db:"id"`
	Name racers.RaceName `db:"name"`
	Date time.Time       `db:"date"`
}

func NewRaces(db *gorm.DB, events Events) Races {
	return Races{db, events, "races"}
}

type Races struct {
	db        *gorm.DB
	events    Events
	tableName string
}

func (r Races) ByID(ctx context.Context, id racers.RaceID) (*racers.Race, error) {
	return nil, nil
}

func (r Races) Save(ctx context.Context, race racers.Race) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.events.Save(ctx, tx, race.ConsumeEvents()...); err != nil {
			return err
		}

		return tx.Table(r.tableName).Save(&raceDB{
			ID:   race.ID(),
			Name: race.Name(),
			Date: time.Time(race.Date()),
		}).Error
	})

	return errors.WrapInternalError(err, "saving race")
}
