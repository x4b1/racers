package test

import (
	"database/sql"
	"errors"
	stdlog "log"
	"os"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/kataras/golog"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/xabi93/racers/internal/api/graph"
	"github.com/xabi93/racers/internal/config"
	"github.com/xabi93/racers/internal/id"
	"github.com/xabi93/racers/internal/log"
	"github.com/xabi93/racers/internal/server"
	"github.com/xabi93/racers/internal/storage/postgres"
	"github.com/xabi93/racers/internal/storage/postgres/test"
)

var conf config.Conf

func TestMain(m *testing.M) {
	var err error
	conf, err = config.Load()
	if err != nil {
		stdlog.Fatalf("loading config tests: %s", err)
	}

	code := initTest(m.Run)

	os.Exit(code)
}

func initTest(f func() int) int {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return failInit("Could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Name:       "racers_db_test",
		Env: []string{
			"POSTGRES_USER=" + conf.Postgres.User,
			"POSTGRES_PASSWORD=" + conf.Postgres.Password,
			"POSTGRES_DB=" + conf.Postgres.Database,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: conf.Postgres.Port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if errors.Is(err, docker.ErrContainerAlreadyExists) {
		err = pool.RemoveContainerByName(opts.Name)
		if err != nil {
			return failInit("Could not start resource: %s", err.Error())
		}
		resource, err = pool.RunWithOptions(&opts)
	}
	if err != nil {
		return failInit("Could not start resource: %s", err.Error())
	}

	resource.Expire(60)
	defer func() {
		panicErr := recover()
		if err := pool.Purge(resource); err != nil {
			golog.Fatalf("Could not purge resource: %s", err)
		}
		if panicErr != nil {
			panic(panicErr)
		}
	}()

	var dbConn *sql.DB
	if err = pool.Retry(func() error {
		dbConn, err = postgres.New(conf.Postgres)
		if err != nil {
			return err
		}
		return dbConn.Ping()
	}); err != nil {
		return failInit("Could not connect to docker: %s", err.Error())
	}
	defer dbConn.Close()

	if err := postgres.RunMigrations(conf.Postgres, dbConn); err != nil {
		return failInit("running migrations %s", err)
	}

	return f()
}

func failInit(format string, v ...interface{}) int {
	stdlog.Printf(format, v...)
	return 1
}

func newSuite(t *testing.T) suite {
	t.Helper()

	db, err := test.New(conf.Postgres)
	if err != nil {
		t.Fatalf("initializing testing conn %s", err)
	}

	srv, err := server.New(conf, log.Noop{}, db)
	if err != nil {
		t.Fatalf("initializing server %s", err)
	}
	return suite{
		graphql: client.New(srv.Handler(), client.Path(server.GraphEndpoint)),
	}
}

type suite struct {
	graphql *client.Client
}

type getRaceResult struct {
	Race struct {
		Typename string `json:"__typename,omitempty"`
		ID       string `json:"id,omitempty"`
		Name     string `json:"name,omitempty"`
		Date     string `json:"date,omitempty"`
		Message  string `json:"message,omitempty"`
	}
}

func getRace(c *client.Client, id id.ID) getRaceResult {
	const query = `query($id: ID!) {
		race(id: $id){
			__typename
			...on Race {
				id
				name
				date
			}
			...on RaceByIDNotFound {
				id
			}
		}
		}`

	var resp getRaceResult

	c.MustPost(query, &resp, client.Var("id", id))

	return resp
}

type createRaceResult struct {
	CreateRace struct {
		Typename string `json:"__typename,omitempty"`
		ID       string `json:"id,omitempty"`
		Name     string `json:"name,omitempty"`
		Date     string `json:"date,omitempty"`
		Message  string `json:"message,omitempty"`
	}
}

func createRace(c *client.Client, req graph.Race) createRaceResult {
	const mutation = `mutation($name: String!, $date: DateTime!) {
		createRace(race:{name:$name,date:$date}){
			__typename
			...on Race {
				id
				name
				date
			}
			...on InvalidNameError {
				message
			}
			...on InvalidDateError {
				message
			}
		}}`

	var resp createRaceResult

	c.MustPost(mutation, &resp,
		client.Var("name", req.Name),
		client.Var("date", req.Date),
	)

	return resp
}
