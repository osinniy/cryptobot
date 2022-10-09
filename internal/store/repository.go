package store

import (
	"osinniy/cryptobot/internal/models"
	"time"
)

type UsersRepository interface {
	// Adds user to storage
	// Returns [ErrNilUser] if user is nil.
	// Returns [ErrUserExists] if user with such ID already exists.
	Add(*models.User) error

	// Checks whether user with such ID exists in storage.
	IsExists(userId int64) (bool, error)

	// Gets user language.
	// Returns [ErrUserNotFound] if user not found.
	Lang(userId int64) (string, error)

	// Sets user language.
	// Returns [ErrUserNotFound] if user not found.
	SetLang(userId int64, lang string) error

	// Returns total stored users count.
	UsersTotal() (int, error)
}

type DataRepository interface {
	// Saves data in storage.
	// Returns [ErrNilData] if provided data is nil.
	// Can return [ErrOldData] in some implementations
	// if provided data is older than last stored.
	Save(*models.Data) error

	// Returns latest saved data from store or error.
	// If there is no data in store, returns both nil.
	Latest() (*models.Data, error)

	// Returns number of stored elements.
	Len() (int, error)

	// Removes entries older than [teardown] time.
	Cleanup(teardown time.Time) (affected int64, err error)
}
