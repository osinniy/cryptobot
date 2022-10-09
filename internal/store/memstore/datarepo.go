package memstore

import (
	"fmt"
	"osinniy/cryptobot/internal/logs"
	"osinniy/cryptobot/internal/models"
	"osinniy/cryptobot/internal/store"
	"sync"
	"time"
)

type DataRepository struct {
	store *MemStore
	data  []models.Data

	mu sync.Mutex
}

func (repo *DataRepository) Save(data *models.Data) error {
	const q = "data.save"

	if data == nil {
		return store.ErrNilData
	}

	repo.mu.Lock()
	if len(repo.data) > 0 {
		if data.Upd <= repo.data[len(repo.data)-1].Upd {
			return store.ErrOldData
		}
	}

	repo.data = append(repo.data, *data)
	repo.mu.Unlock()

	logs.MemQuery(&repo.store.logger, q, nil, map[string]any{"data": data})
	return nil
}

func (repo *DataRepository) Latest() (data *models.Data, err error) {
	const q = "data.latest"

	repo.mu.Lock()
	if len(repo.data) > 0 {
		data = &repo.data[len(repo.data)-1]
	}
	repo.mu.Unlock()

	logs.MemQuery(&repo.store.logger, q, err, map[string]any{"data": data})
	return
}

func (repo *DataRepository) Len() (int, error) {
	const q = "data.len"

	repo.mu.Lock()
	len := len(repo.data)
	repo.mu.Unlock()

	logs.MemQuery(&repo.store.logger, q, nil, map[string]any{"len": len})
	return len, nil
}

func (repo *DataRepository) Cleanup(teardown time.Time) (affected int64, err error) {
	const q = "data.cleanup"
	unix := teardown.Unix()
	dur := time.Since(teardown)

	repo.mu.Lock()

	len := len(repo.data)
	var found bool
	// Find first satisfying element and remove all elements before it
	for i, data := range repo.data {
		if data.Upd >= unix {
			repo.data = repo.data[i:]
			affected = int64(len - (len - i))
			found = true
			break
		}
	}
	// If we can't find any satisfying elements, clear the whole slice
	if !found {
		repo.data = repo.data[:0]
		affected = int64(len)
	}

	repo.mu.Unlock()

	logs.MemQuery(&repo.store.logger, q, nil, map[string]any{
		"affected": affected,
		"teardown": fmt.Sprintf("%1.fh", dur.Hours()),
	})
	return
}
