package postgres

import (
	"context"
	"time"

	racers "github.com/xabi93/racers/internal"
	"github.com/xabi93/racers/internal/service"
	"gorm.io/gorm"
)

type race struct {
	ID      racers.RaceID   `db:"id"`
	Name    racers.RaceName `db:"name"`
	Date    time.Time       `db:"date"`
	OwnerID racers.UserID   `db:"owner_id"`
}

func (race) TableName() string {
	return "races"
}

func (r race) toDomain() racers.Race {
	return racers.Race{
		ID:    r.ID,
		Name:  r.Name,
		Date:  racers.RaceDate(r.Date),
		Owner: r.OwnerID,
	}
}

func NewRaces(db *gorm.DB) Races {
	return Races{Repository{db}}
}

type Races struct {
	repo Repository
}

func (r Races) All(ctx context.Context) ([]racers.Race, error) {
	db := r.repo.DB(ctx)
	rows, err := db.Model(&race{}).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]racers.Race, 0)
	for rows.Next() {
		var dbRace race
		if err := db.ScanRows(rows, &dbRace); err != nil {
			return nil, err
		}
		result = append(result, dbRace.toDomain())
	}

	return result, nil
}

func (r Races) Get(ctx context.Context, id racers.RaceID) (racers.Race, error) {
	var raceDB race
	if err := r.repo.DB(ctx).Take(&raceDB, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return racers.Race{}, service.ErrRaceNotFound
		}
		return racers.Race{}, err
	}

	return raceDB.toDomain(), nil
}

func (r Races) Exists(ctx context.Context, in racers.Race) (bool, error) {
	var count int64
	query := r.repo.DB(ctx).
		Model(&race{}).
		Where(&race{Name: in.Name, Date: time.Time(in.Date)})
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r Races) Save(ctx context.Context, in racers.Race) error {
	return r.repo.DB(ctx).Create(&race{
		ID:      in.ID,
		Name:    in.Name,
		Date:    time.Time(in.Date),
		OwnerID: in.Owner,
	}).Error
}
