package sqlstore

import (
	sqlLib "database/sql"
	"osinniy/cryptobot/internal/logs"
	"osinniy/cryptobot/internal/models"
	"osinniy/cryptobot/internal/store"
	"sync"
	"time"
)

type UsersRepository struct {
	store *SqlStore
	mu    sync.RWMutex
}

// Inserts recipient id, current time and language code into users table.
func (repo *UsersRepository) Add(user *models.User) (err error) {
	const q = "users.add"
	const sql = `INSERT INTO users (id, first_seen, lang) VALUES ($1, $2, $3)`

	if user == nil {
		return store.ErrNilUser
	}

	timestamp := time.Now()
	repo.mu.Lock()
	result, err := repo.store.db.Exec(sql, user.Id, user.FirstSeen, user.Lang)
	repo.mu.Unlock()

	logs.DbQuery(&repo.store.logger, q, err, timestamp, map[string]any{"user": user})
	if err != nil {
		return err
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAff == 0 {
		err = store.ErrUserExists
	}
	return
}

func (repo *UsersRepository) IsExists(userId int64) (exists bool, err error) {
	const q = "users.is_exists"
	const sql = `SELECT EXISTS(SELECT 1 FROM users WHERE id=?)`

	timestamp := time.Now()
	repo.mu.RLock()
	err = repo.store.db.Get(&exists, sql, userId)
	repo.mu.RUnlock()

	logs.DbQuery(&repo.store.logger, q, err, timestamp, map[string]any{"user_id": userId, "exists": exists})
	return
}

func (repo *UsersRepository) Lang(userId int64) (lang string, err error) {
	const q = "users.lang"
	const sql = `SELECT lang FROM users WHERE id=?`

	timestamp := time.Now()
	repo.mu.RLock()
	err = repo.store.db.Get(&lang, sql, userId)
	repo.mu.RUnlock()

	if err == sqlLib.ErrNoRows {
		err = store.ErrUserNotFound
	}
	logs.DbQuery(&repo.store.logger, q, err, timestamp, map[string]any{"user_id": userId, "lang": lang})
	return
}

func (repo *UsersRepository) SetLang(userId int64, lang string) (err error) {
	const q = "users.set_lang"
	const sql = `UPDATE users SET lang=? WHERE id=?`

	timestamp := time.Now()
	repo.mu.Lock()
	result, err := repo.store.db.Exec(sql, lang, userId)
	repo.mu.Unlock()

	logs.DbQuery(&repo.store.logger, q, err, timestamp, map[string]any{"user_id": userId, "lang": lang})
	if err != nil {
		return err
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAff == 0 {
		err = store.ErrUserNotFound
	}
	return
}

func (repo *UsersRepository) UsersTotal() (total int, err error) {
	const q = "users.total"
	const sql = `SELECT COUNT(*) FROM users`

	timestamp := time.Now()
	repo.mu.RLock()
	err = repo.store.db.Get(&total, sql)
	repo.mu.RUnlock()

	logs.DbQuery(&repo.store.logger, q, err, timestamp, map[string]any{"total": total})
	return
}
