package memstore_test

import (
	"os"
	"osinniy/cryptobot/internal/store/memstore"
	"reflect"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.Disabled)

	os.Exit(m.Run())
}

func TestMemStore_Open(t *testing.T) {
	t.Parallel()

	store := memstore.New()
	err := store.Open()
	assert.NoError(t, err)
}

func TestMemStore_Close(t *testing.T) {
	t.Parallel()

	store := memstore.New()
	err := store.Open()
	assert.NoError(t, err)
	store.Users()
	store.Data()

	err = store.Close()
	assert.NoError(t, err)

	usersRepo := reflect.ValueOf(store).Elem().FieldByName("usersRepo")
	assert.NotEqualValues(t, usersRepo, reflect.Value{}, "MemStore.usersRepo field not found")
	dataRepo := reflect.ValueOf(store).Elem().FieldByName("dataRepo")
	assert.NotEqualValues(t, dataRepo, reflect.Value{}, "MemStore.dataRepo field not found")

	assert.True(t, usersRepo.IsNil(), "MemStore.usersRepo field is not nil")
	assert.True(t, dataRepo.IsNil(), "MemStore.dataRepo field is not nil")
}
