package sqlstore_test

import (
	"osinniy/cryptobot/internal/models"
	"osinniy/cryptobot/internal/store"
	"osinniy/cryptobot/internal/store/sqlstore"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

const USERS_TABLE = "users"

func TestUsersRepo_Add(t *testing.T) {
	s, teardown := sqlstore.TestStore(t)
	defer teardown(USERS_TABLE)

	testCases := []struct {
		name string
		args *models.User
		exp  error
	}{
		{
			name: "ok",
			args: models.TestUser(t),
			exp:  nil,
		},
		{
			name: "duplicate",
			args: models.TestUser(t),
			exp:  store.ErrUserExists,
		},
		{
			name: "nil user",
			args: nil,
			exp:  store.ErrNilUser,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.Users().Add(tc.args)
			assert.Equal(t, tc.exp, err)
		})
	}
}

func TestUsersRepo_IsExists(t *testing.T) {
	s, teardown := sqlstore.TestStore(t)
	defer teardown(USERS_TABLE)

	user := models.TestUser(t)

	t.Run("not_exists", func(t *testing.T) {
		exists, err := s.Users().IsExists(user.Id)
		assert.False(t, exists)
		assert.NoError(t, err)
	})

	t.Run("exists", func(t *testing.T) {
		err := s.Users().Add(user)
		assert.NoError(t, err)

		exists, err := s.Users().IsExists(user.Id)
		assert.True(t, exists)
		assert.NoError(t, err)
	})
}

func TestUsersRepo_Lang(t *testing.T) {
	s, teardown := sqlstore.TestStore(t)
	defer teardown(USERS_TABLE)

	t.Run("normal", func(t *testing.T) {
		err := s.Users().Add(models.TestUser(t))
		assert.NoError(t, err)

		lang, err := s.Users().Lang(1)
		assert.NoError(t, err)
		assert.Equal(t, "en", lang)
	})

	t.Run("user_not_found", func(t *testing.T) {
		lang, err := s.Users().Lang(-1)
		assert.Equal(t, "", lang)
		assert.ErrorIs(t, err, store.ErrUserNotFound)
	})
}

func TestUsersRepo_SetLang(t *testing.T) {
	s, teardown := sqlstore.TestStore(t)
	defer teardown(USERS_TABLE)

	user := models.TestUser(t)

	t.Run("normal", func(t *testing.T) {
		err := s.Users().Add(user)
		assert.NoError(t, err)

		err = s.Users().SetLang(1, "ua")
		assert.NoError(t, err)

		lang, err := s.Users().Lang(1)
		assert.NoError(t, err)
		assert.Equal(t, "ua", lang)
	})

	t.Run("user_not_found", func(t *testing.T) {
		err := s.Users().SetLang(-1, "ua")
		assert.ErrorIs(t, err, store.ErrUserNotFound)
	})
}

func TestUsersRepo_UsersTotal(t *testing.T) {
	s, teardown := sqlstore.TestStore(t)
	defer teardown(USERS_TABLE)

	t.Run("zero", func(t *testing.T) {
		total, err := s.Users().UsersTotal()
		assert.NoError(t, err)
		assert.Equal(t, 0, total)
	})

	t.Run("one", func(t *testing.T) {
		err := s.Users().Add(models.TestUser(t))
		assert.NoError(t, err)

		total, err := s.Users().UsersTotal()
		assert.NoError(t, err)
		assert.Equal(t, 1, total)
	})
}

func BenchmarkUsersRepo(b *testing.B) {
	b.ReportAllocs()

	store, teardown := sqlstore.TestStore(b)
	defer teardown(USERS_TABLE)

	repo := store.Users()

	users := make([]models.User, 0, b.N)
	for i := 0; i < b.N; i++ {
		users = append(users, models.User{
			Id:        int64(i),
			FirstSeen: int64(i),
			Lang:      "en",
		})
	}

	zerolog.SetGlobalLevel(zerolog.Disabled)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.Add(&users[i])
		repo.IsExists(int64(i))
		repo.SetLang(users[i].Id, "ua")
		repo.Lang(users[i].Id)
		repo.UsersTotal()
	}
}
