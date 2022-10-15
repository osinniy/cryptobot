package store

// Store is a common interface for all stores.
// This may be sql, memory, and so on.
// Store is safe for concurrent use.
type Store interface {
	// Opens storage so you can obtain repositories from it.
	// This might not be required for some stores.
	Open() error

	// Closes storage and nullifies repositories so they will be recreated
	// on next access if you want to open storage again later.
	Close() error

	// Returns [UsersRepository]
	Users() UsersRepository

	// Returns [DataRepository]
	Data() DataRepository
}
