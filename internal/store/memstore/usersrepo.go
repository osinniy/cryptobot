package memstore

import (
	"osinniy/cryptobot/internal/logs"
	"osinniy/cryptobot/internal/models"
	"osinniy/cryptobot/internal/store"
	"sync"
)

type UsersRepository struct {
	store *MemStore
	data  map[int64]*models.User

	mu sync.RWMutex
}

func (repo *UsersRepository) Add(user *models.User) (err error) {
	const q = "users.add"

	if user == nil {
		return store.ErrNilUser
	}

	repo.mu.RLock()
	_, ok := repo.data[user.Id]
	repo.mu.RUnlock()

	if ok {
		return store.ErrUserExists
	}

	repo.mu.Lock()
	repo.data[user.Id] = user
	repo.mu.Unlock()

	logs.MemQuery(&repo.store.logger, q, err, map[string]any{"user": user})
	return
}

func (repo *UsersRepository) IsExists(userId int64) (bool, error) {
	const q = "users.is_exists"

	repo.mu.RLock()
	_, ok := repo.data[userId]
	repo.mu.RUnlock()

	logs.MemQuery(&repo.store.logger, q, nil, map[string]any{"user_id": userId, "exists": ok})
	return ok, nil
}

func (repo *UsersRepository) Lang(userId int64) (lang string, err error) {
	const q = "users.lang"

	repo.mu.RLock()
	user, ok := repo.data[userId]
	repo.mu.RUnlock()

	if !ok {
		err = store.ErrUserNotFound
	}
	if user != nil {
		lang = user.Lang
	}

	logs.MemQuery(&repo.store.logger, q, err, map[string]any{"user_id": userId, "lang": lang})
	return lang, err
}

func (repo *UsersRepository) SetLang(userId int64, lang string) (err error) {
	const q = "users.set_lang"

	repo.mu.RLock()
	user, ok := repo.data[userId]
	repo.mu.RUnlock()

	if !ok {
		err = store.ErrUserNotFound
	}
	if user != nil {
		user.Lang = lang
	}

	logs.MemQuery(&repo.store.logger, q, err, map[string]any{"user_id": userId, "lang": lang})
	return
}

func (repo *UsersRepository) UsersTotal() (total int, err error) {
	const q = "users.total"

	repo.mu.Lock()
	total = len(repo.data)
	repo.mu.Unlock()

	logs.MemQuery(&repo.store.logger, q, err, map[string]any{"total": total})
	return
}
