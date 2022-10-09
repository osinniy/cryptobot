package mergedstore

import (
	"osinniy/cryptobot/internal/models"
	"osinniy/cryptobot/internal/store"
	"osinniy/cryptobot/internal/store/memstore"
	"osinniy/cryptobot/internal/store/sqlstore"
)

type MergedUsersRepository struct {
	mem *memstore.UsersRepository
	db  *sqlstore.UsersRepository
}

func (repo *MergedUsersRepository) Add(user *models.User) error {
	if err := repo.mem.Add(user); err != nil && err != store.ErrUserExists {
		return err
	}
	if err := repo.db.Add(user); err != nil {
		return err
	}

	return nil
}

func (repo *MergedUsersRepository) IsExists(userId int64) (bool, error) {
	exists, err := repo.mem.IsExists(userId)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}

	return repo.db.IsExists(userId)
}

func (repo *MergedUsersRepository) Lang(userId int64) (string, error) {
	lang, err := repo.mem.Lang(userId)
	if err != nil && err != store.ErrUserNotFound {
		return "", err
	}
	if lang != "" {
		return lang, nil
	}

	return repo.db.Lang(userId)
}

func (repo *MergedUsersRepository) SetLang(userId int64, lang string) error {
	if err := repo.mem.SetLang(userId, lang); err != nil && err != store.ErrUserNotFound {
		return err
	}
	return repo.db.SetLang(userId, lang)
}

// Always returns number from database as it is more accurate
func (repo *MergedUsersRepository) UsersTotal() (int, error) {
	return repo.db.UsersTotal()
}
