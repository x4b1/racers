package postgres

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
