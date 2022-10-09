package sqlstore

import (
	"osinniy/cryptobot/internal/store"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	MODULE = "store"

	MAX_IDLE_CONS = 20
)

type SqlStore struct {
	Path string

	db     *sqlx.DB
	logger zerolog.Logger

	usersRepo *UsersRepository
	dataRepo  *DataRepository
}

// Initializes and opens database.
// Aborts if database can't be initialized
func Init(path string) *SqlStore {
	store := New(path)

	err := store.Open()
	if err != nil {
		store.logger.Fatal().Err(err).Str("db", path).Msg("failed to connect to database")
	} else {
		store.logger.Info().Str("db", path).Msg("database connected")
	}

	return store
}

// Creates empty storage.
// Do not forget to open it
func New(path string) *SqlStore {
	return &SqlStore{
		Path:   path,
		logger: log.With().Str("module", MODULE).Logger(),
	}
}

// Connect to a database and verify with a ping
func (s *SqlStore) Open() error {
	db, err := sqlx.Connect("sqlite3", s.Path)
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(MAX_IDLE_CONS)

	s.db = db
	return nil
}

func (s *SqlStore) Close() error {
	if s.db != nil {
		if s.usersRepo != nil {
			s.usersRepo = nil
		}
		if s.dataRepo != nil {
			s.dataRepo = nil
		}
		return s.db.Close()
	}
	return nil
}

func (s *SqlStore) Users() store.UsersRepository {
	if s.usersRepo == nil {
		s.usersRepo = &UsersRepository{
			store: s,
		}
	}

	return s.usersRepo
}

func (s *SqlStore) Data() store.DataRepository {
	if s.dataRepo == nil {
		s.dataRepo = &DataRepository{
			store: s,
		}
	}

	return s.dataRepo
}
