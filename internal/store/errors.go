package store

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrNilUser      = errors.New("user is nil")
	ErrUserExists   = errors.New("user already exists")

	ErrNilData = errors.New("data is nil")
	ErrOldData = errors.New("data is older than latest saved data")
)
