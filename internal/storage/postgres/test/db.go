package test

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/xabi93/racers/internal/storage/postgres"
)

const testDriver = "test-conn"

var once sync.Once

func New(c postgres.Config) (*sql.DB, error) {
	once.Do(func() {
		txdb.Register(testDriver, "postgres", postgres.Url(c))
	})
	return sql.Open(testDriver, fmt.Sprintf("connection_%d", time.Now().UnixNano()))
}
