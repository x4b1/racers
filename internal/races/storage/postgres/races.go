package postgres

import (
	"context"
	"time"

	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/errors"
	"github.com/xabi93/racers/internal/races"
	"gorm.io/gorm"
)

type raceDB struct {
	ID   races.RaceID   `db:"id"`
	Name races.RaceName `db:"name"`
	Date time.Time      `db:"date"`
}

func NewRaces(db *gorm.DB, events Events) Races {
	return Races{db, events, "races"}
}

type Races struct {
	db        *gorm.DB
	events    Events
	tableName string
}

func (r Races) ByID(ctx context.Context, id races.RaceID) (*racers.Race, error) {
	return nil, nil
}

func (r Races) Save(ctx context.Context, race races.Race) error {
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

	return errors.WrapInternalError(err)
}
