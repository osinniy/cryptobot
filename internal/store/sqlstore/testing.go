package sqlstore

import (
	"osinniy/cryptobot/internal/store"
	"testing"

	"github.com/jmoiron/sqlx"
)

const (
	ROOT_PATH = "../../../"
	DB_PATH   = "test.sqlite"
)

func TestStore(t testing.TB, testDb ...string) (store.Store, func(...string)) {
	t.Helper()
	if len(testDb) == 0 || testDb[0] == "" {
		testDb = []string{DB_PATH}
	}

	s, _, err := TestDb(t, testDb[0])
	if err != nil {
		t.Fatal(err)
	}

	return s, func(tables ...string) {
		if len(tables) > 0 {
			for _, table := range tables {
				s.db.MustExec("DELETE FROM " + table)
			}
		}

		if err := s.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDb(t testing.TB, testDb string) (*SqlStore, *sqlx.DB, error) {
	t.Helper()

	s := New(testDb)
	err := s.Open()
	if err != nil {
		return s, nil, err
	}
	return s, s.db, nil
}
