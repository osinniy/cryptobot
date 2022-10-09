package sqlstore

import (
	sqlLib "database/sql"
	"fmt"
	"osinniy/cryptobot/internal/logs"
	"osinniy/cryptobot/internal/models"
	"osinniy/cryptobot/internal/store"
	"sync"
	"time"

	"github.com/mattn/go-sqlite3"
)

type DataRepository struct {
	store *SqlStore
	mu    sync.RWMutex
}

func (repo *DataRepository) Save(data *models.Data) (err error) {
	const q = "data.save"
	const sql = `INSERT INTO data (
		upd, btcd, ethd, btcd_ch, ethd_ch, stablecoinscap,
		totalcap, totalcap_ch, liquid_usd, liquid_usd_ch, liquid_num, oi, oi_ch
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	if data == nil {
		return store.ErrNilData
	}

	timestamp := time.Now()
	repo.mu.Lock()
	_, err = repo.store.db.Exec(sql,
		data.Upd,
		data.BTCDom,
		data.ETHDom,
		data.BTCDom24HChange,
		data.ETHDom24HChange,
		data.StablecoinsCap,
		data.TotalCap,
		data.TotalCap24HChange,
		data.Liquidation24HUsd,
		data.Liquidations24HUsdChange,
		data.Liquidation24HNum,
		data.OpenInterest,
		data.OpenInterest24HChange,
	)
	repo.mu.Unlock()

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.Code == 19 {
				err = store.ErrOldData
			}
		}
	}

	logs.DbQuery(&repo.store.logger, q, err, timestamp, map[string]any{"data": data})

	return
}

func (repo *DataRepository) Latest() (data *models.Data, err error) {
	const q = "data.latest"
	const sql = "SELECT * FROM data ORDER BY upd DESC LIMIT 1"

	timestamp := time.Now()
	repo.mu.RLock()
	row := repo.store.db.QueryRowx(sql)
	repo.mu.RUnlock()

	data = &models.Data{}
	err = row.Scan(
		&data.Upd,
		&data.BTCDom,
		&data.ETHDom,
		&data.BTCDom24HChange,
		&data.ETHDom24HChange,
		&data.StablecoinsCap,
		&data.TotalCap,
		&data.TotalCap24HChange,
		&data.Liquidation24HUsd,
		&data.Liquidations24HUsdChange,
		&data.Liquidation24HNum,
		&data.OpenInterest,
		&data.OpenInterest24HChange,
	)
	if err == sqlLib.ErrNoRows {
		data = nil
		err = nil
	}

	logs.DbQuery(&repo.store.logger, q, err, timestamp, map[string]any{"data": data})

	return
}

func (repo *DataRepository) Len() (len int, err error) {
	const q = "data.len"
	const sql = "SELECT COUNT(*) FROM data"

	timestamp := time.Now()
	repo.mu.RLock()
	err = repo.store.db.Get(&len, sql)
	repo.mu.RUnlock()

	logs.DbQuery(&repo.store.logger, q, err, timestamp, map[string]any{"len": len})
	return
}

func (repo *DataRepository) Cleanup(teardown time.Time) (affected int64, err error) {
	const q = "data.cleanup"
	const sql = "DELETE FROM data WHERE upd < ?"

	dur := time.Since(teardown)
	timestamp := time.Now()

	repo.mu.Lock()
	result, err := repo.store.db.Exec(sql, teardown.Unix())
	repo.mu.Unlock()
	if err != nil {
		return
	}

	affected, err = result.RowsAffected()
	if err != nil {
		return
	}

	logs.DbQuery(&repo.store.logger, q, err, timestamp, map[string]any{
		"affected": affected,
		"teardown": fmt.Sprintf("%1.fh", dur.Hours()),
	})

	return
}
