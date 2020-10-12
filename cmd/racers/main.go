package main

import (
	"fmt"
	"io"
	"os"

	"github.com/kr/pretty"
	"github.com/xabi93/racers/internal/config"
	"github.com/xabi93/racers/internal/log"
	"github.com/xabi93/racers/internal/server"
	"github.com/xabi93/racers/internal/storage/postgres"
)

func main() {
	if err := Run(os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func Run(out io.Writer) error {
	conf, err := config.Load()
	if err != nil {
		return err
	}

	log, err := log.NewZap(conf.Env, out)
	if err != nil {
		return err
	}

	db, err := postgres.New(conf.Postgres)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := postgres.RunMigrations(conf.Postgres, db); err != nil {
		pretty.Println(err)
		return err
	}

	s, err := server.New(conf, log, db)
	if err != nil {
		return err
	}

	return s.Serve()
}
