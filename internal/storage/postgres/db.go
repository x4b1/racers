package postgres

import (
	"database/sql"
	"fmt"
	"net"
	"net/url"

	_ "github.com/lib/pq"
)

func New(c Config) (*sql.DB, error) {
	conn, err := sql.Open("postgres", Url(c))
	if err != nil {
		return nil, err
	}

	return conn, conn.Ping()
}

func Url(c Config) string {
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
