package server

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"

	"github.com/xabi93/racers/internal/instrumentation/log"
	"github.com/xabi93/racers/internal/server/graph"
	"github.com/xabi93/racers/internal/server/graph/instrumentation"
	"github.com/xabi93/racers/internal/service"
	"github.com/xabi93/racers/internal/storage/postgres"
	"github.com/xabi93/racers/internal/users"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const GraphEndpoint = "/graph"

func New(conf Conf, logger log.Logger, db *sql.DB, uProvider users.UsersProvider) (Server, error) {
	s := Server{
		conf:   conf,
		logger: logger,
		db:     db,
		users:  users.Users{UsersProvider: uProvider},
	}

	if err := s.initService(); err != nil {
		return Server{}, err
	}

	s.initHandler()

	return s, nil
}

type Server struct {
	conf   Conf
	logger log.Logger
	db     *sql.DB
	users  users.Users

	handler http.Handler

	races service.Races
}

func (s *Server) initService() error {
	db, err := postgres.New(s.db)
	if err != nil {
		return err
	}

	eventsRepo := postgres.NewEvents(db)
	racesRepo := postgres.NewRaces(db)

	s.races = service.NewRaces(racesRepo, s.users, postgres.TransactionFactory(db), eventsRepo)

	return nil
}

func (s *Server) initHandler() {
	r := mux.NewRouter()

	r.Use(users.AuthMiddleware(s.users))

	r.Handle("/playground", playground.Handler("racers", GraphEndpoint))

	graphServer := handler.NewDefaultServer(graph.NewExecutableSchema(graph.New(s.races)))

	promRegistry := prometheus.NewRegistry()
	graphServer.Use(instrumentation.NewPrometheus(promRegistry, "racers"))
	r.Handle(GraphEndpoint, graphServer)

	r.Handle("/metrics", promhttp.InstrumentMetricHandler(
		promRegistry, promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}),
	))

	s.handler = r
}

func (s *Server) Handler() http.Handler {
	return s.handler
}

func (s *Server) Serve() error {
	addr := net.JoinHostPort("", s.conf.Port)

	s.logger.Info(context.Background(), fmt.Sprintf("Server running on: %s", addr), nil)

	return http.ListenAndServe(addr, s.handler)
}
