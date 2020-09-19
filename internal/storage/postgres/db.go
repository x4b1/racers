package postgres

import (
	"database/sql"
	"fmt"
	"net"
	"net/url"

	_ "github.com/lib/pq"
)

type Config struct {
	User         string `env:"DATABASE_USER" envDefault:"racers"`
	Password     string `env:"DATABASE_PASS" envDefault:"racers"`
	Host         string `env:"DATABASE_HOST" envDefault:"localhost"`
	Port         string `env:"DATABASE_PORT" envDefault:"5433"`
	Database     string `env:"DATABASE_NAME" envDefault:"racers"`
	SSLMode      string `env:"DATABASE_SSL_MODE" envDefault:"disable"`
	BinaryParams string `env:"DATABASE_BINARY_PARAMS" envDefault:"yes"`

	MigrationsTable string `env:"DATABASE_MIGRATIONS_TABLE" envDefault:"migrations"`
}

func (c Config) URL() string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, c.Port),
		Path:     c.Database,
		RawQuery: fmt.Sprintf("sslmode=%s&binary_parameters=%s", c.SSLMode, c.BinaryParams),
	}
	u.Query().Add("sslmode", c.SSLMode)
	u.Query().Add("binary_parameters", c.BinaryParams)

	return u.String()
}

func Connect(c Config) (*sql.DB, error) {
	conn, err := sql.Open("postgres", c.URL())
	if err != nil {
		return nil, err
	}

	return conn, conn.Ping()
}
