package gorm

import (
	"database/sql"

	"github.com/xabi93/racers/internal/errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(conn *sql.DB) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: conn}), &gorm.Config{})

	return db, errors.WrapInternalError(err, "initializing gorm client")
}
