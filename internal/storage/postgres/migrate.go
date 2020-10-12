package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(conf Config, db *sql.DB) error {
	dbInstance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationsPath()), conf.Database, dbInstance)
	if err != nil {
		return err
	}

	err = m.Up()
	if err == nil || err == migrate.ErrNoChange {
		return nil
	}

	return err
}

func migrationsPath() string {
	if _, err := os.Stat("migrations"); err == nil {
		return "migrations"
	}

	projectRoot, _ := os.Getwd()
	for !strings.HasSuffix(projectRoot, "/racers") {
		projectRoot = filepath.Dir(projectRoot)
	}

	return projectRoot + "/internal/storage/postgres/migrations"
}
