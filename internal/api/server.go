package api

import (
	"net/http"
	"osinniy/cryptobot/internal/store"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	MODULE       = "api"
	DEFAULT_PORT = 1488
)

type Server struct {
	port   int
	router *mux.Router
	store  store.Store
	logger zerolog.Logger
}

// Creates new api server with given port.
// If port is not specified, uses DEFAULT_PORT.
// You should not disclose this port to the public internet
// because it contains sensitive data.
func NewServer(store store.Store, port ...int) *Server {
	if len(port) == 0 || port[0] == 0 {
		port = []int{DEFAULT_PORT}
	}

	server := &Server{
		port:   port[0],
		router: mux.NewRouter().StrictSlash(true),
		store:  store,
		logger: log.With().Str("module", MODULE).Logger(),
	}

	server.configureRoutes()
	server.configureMIddlewares()

	return server
}

func (s *Server) configureRoutes() {
	s.router.HandleFunc("/", s.root)
	s.router.HandleFunc("/users/total", s.usersTotal)
	s.router.HandleFunc("/users/lang", s.usersLang)
}

func (s *Server) configureMIddlewares() {
	s.router.Use(s.mLogger)
}

func (s *Server) Serve() {
	s.logger.Info().Int("port", s.port).Msg("starting api server")

	srv := http.Server{
		Addr:         ":" + strconv.Itoa(s.port),
		Handler:      s.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		s.logger.Error().Err(err).Int("port", s.port).Msg("unable to start api server")
	}
}
