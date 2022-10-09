package memstore

import (
	"fmt"
	"osinniy/cryptobot/internal/models"
	"osinniy/cryptobot/internal/store"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	MODULE = "store"

	DATAREP_SLICE_PREALLOC = 1024
)

type MemStore struct {
	logger zerolog.Logger

	usersRepo *UsersRepository
	dataRepo  *DataRepository
}

// Creates new mem store ready to use
func New() *MemStore {
	return &MemStore{
		logger: log.With().Str("module", MODULE).Logger(),
	}
}

func (s *MemStore) Open() error {
	return nil
}

func (s *MemStore) Close() error {
	s.usersRepo = nil
	s.dataRepo = nil
	return nil
}

func (s *MemStore) Users() store.UsersRepository {
	if s.usersRepo == nil {
		s.usersRepo = &UsersRepository{
			store: s,
			data:  make(map[int64]*models.User),
		}
		fmt.Println(s.usersRepo.data)
	}

	return s.usersRepo
}

func (s *MemStore) Data() store.DataRepository {
	if s.dataRepo == nil {
		s.dataRepo = &DataRepository{
			store: s,
			data:  make([]models.Data, 0, DATAREP_SLICE_PREALLOC),
		}
	}

	return s.dataRepo
}
