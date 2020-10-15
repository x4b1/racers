package server

// const GraphEndpoint = "/graph"

// func New(conf config.Conf, logger log.Logger, db *sql.DB) (Server, error) {
// 	s := Server{conf: conf, logger: logger, db: db}

// 	if err := s.initService(); err != nil {
// 		return Server{}, err
// 	}

// 	s.initHandler()

// 	return s, nil
// }

// type Server struct {
// 	conf   config.Conf
// 	logger log.Logger
// 	db     *sql.DB

// 	handler http.Handler

// 	service *service.Service
// }

// func (s *Server) initService() error {
// 	db, err := gorm.New(s.db)
// 	if err != nil {
// 		return err
// 	}

// 	eventsRepo := gorm.NewEvents(db)
// 	racesRepo := gorm.NewRaces(db, eventsRepo)

// 	s.service = service.NewService(
// 		service.NewRacesMetrics(service.NewRacesService(racesRepo, nil)),
// 		nil,
// 	)

// 	return nil
// }

// func (s *Server) initHandler() {
// 	r := mux.NewRouter()

// 	r.Handle("/playground", playground.Handler("racers", GraphEndpoint))
// 	r.Handle(GraphEndpoint, handler.NewDefaultServer(graph.NewExecutableSchema(graph.New(s.service, ent.New(s.conf.Postgres, s.logger, s.db)))))

// 	s.handler = r
// }

// func (s *Server) Handler() http.Handler {
// 	return s.handler
// }

// func (s *Server) Serve() error {
// 	addr := net.JoinHostPort("", s.conf.Port)

// 	s.logger.Info(context.Background(), fmt.Sprintf("Server running on: %s", addr), nil)

// 	return http.ListenAndServe(addr, s.handler)
// }
