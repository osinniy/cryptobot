package mergedstore

import (
	"osinniy/cryptobot/internal/models"
	"osinniy/cryptobot/internal/store/memstore"
	"osinniy/cryptobot/internal/store/sqlstore"
	"time"
)

type MergedDataRepository struct {
	mem *memstore.DataRepository
	db  *sqlstore.DataRepository
}

func (repo *MergedDataRepository) Save(data *models.Data) error {
	if err := repo.mem.Save(data); err != nil {
		return err
	}
	if err := repo.db.Save(data); err != nil {
		return err
	}

	return nil
}

func (repo *MergedDataRepository) Latest() (*models.Data, error) {
	data, err := repo.mem.Latest()
	if err != nil {
		return nil, err
	}
	if data != nil {
		return data, nil
	}

	// If there is no data in memory, try to get it from database
	data, err = repo.db.Latest()
	if err != nil {
		return nil, err
	}
	if data != nil {
		// Save data to memory if it was found in database
		if err := repo.mem.Save(data); err != nil {
			return nil, err
		}
	}

	return data, nil
}

// Always returns number from database as it is more accurate
func (repo *MergedDataRepository) Len() (int, error) {
	return repo.db.Len()
}

func (repo *MergedDataRepository) Cleanup(teardown time.Time) (affected int64, err error) {
	if affected, err = repo.mem.Cleanup(teardown); err != nil {
		return
	}
	if affected, err = repo.db.Cleanup(teardown); err != nil {
		return
	}

	return
}
