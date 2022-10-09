package mergedstore

import (
	"osinniy/cryptobot/internal/store"
	"osinniy/cryptobot/internal/store/memstore"
	"osinniy/cryptobot/internal/store/sqlstore"
)

// MergedStore stores data in database and uses memory as cache.
type MergedStore struct {
	mem *memstore.MemStore
	db  *sqlstore.SqlStore

	usersRepo *MergedUsersRepository
	dataRepo  *MergedDataRepository
}

// Combines memory and database stores into one.
// You mustn't use provided memory and database stores separately after creating MergedStore.
func Merge(mem *memstore.MemStore, db *sqlstore.SqlStore) *MergedStore {
	return &MergedStore{
		mem: mem,
		db:  db,
	}
}

func (s *MergedStore) Open() error {
	if err := s.mem.Open(); err != nil {
		return err
	}
	if err := s.db.Open(); err != nil {
		return err
	}

	return nil
}

func (s *MergedStore) Close() error {
	if err := s.mem.Close(); err != nil {
		return err
	}
	if err := s.db.Close(); err != nil {
		return err
	}

	return nil
}

func (s *MergedStore) Users() store.UsersRepository {
	if s.usersRepo == nil {
		s.usersRepo = &MergedUsersRepository{
			mem: s.mem.Users().(*memstore.UsersRepository),
			db:  s.db.Users().(*sqlstore.UsersRepository),
		}
	}

	return s.usersRepo
}

func (s *MergedStore) Data() store.DataRepository {
	if s.dataRepo == nil {
		s.dataRepo = &MergedDataRepository{
			mem: s.mem.Data().(*memstore.DataRepository),
			db:  s.db.Data().(*sqlstore.DataRepository),
		}
	}

	return s.dataRepo
}
