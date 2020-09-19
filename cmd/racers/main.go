package main

import (
	"fmt"
	"io"
	"os"

	"github.com/xabi93/racers/internal/instrumentation/log"
	"github.com/xabi93/racers/internal/server"
	"github.com/xabi93/racers/internal/storage/postgres"
	"github.com/xabi93/racers/internal/users"
)

func main() {
	if err := Run(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func Run(out io.Writer) error {
	conf, err := server.LoadConf()
	if err != nil {
		return err
	}

	log, err := log.NewLogrus(out)
	if err != nil {
		return err
	}

	db, err := postgres.Connect(conf.Postgres)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := postgres.RunMigrations(conf.Postgres, db); err != nil {
		return err
	}

	s, err := server.New(conf, log, db, users.Mock{})
	if err != nil {
		return err
	}

	return s.Serve()
}
