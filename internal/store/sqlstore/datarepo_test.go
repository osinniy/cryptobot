package sqlstore_test

import (
	"osinniy/cryptobot/internal/models"
	"osinniy/cryptobot/internal/store"
	"osinniy/cryptobot/internal/store/sqlstore"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

const DATA_TABLE = "data"

func TestDataRepo_Save(t *testing.T) {
	testCases := []struct {
		name    string
		prepare func() []*models.Data
		arg     func() *models.Data
		exp     error
	}{
		{
			name: "ok",
			arg: func() *models.Data {
				return models.TestData(t)
			},
			exp: nil,
		},
		{
			name: "nil data",
			arg: func() *models.Data {
				return nil
			},
			exp: store.ErrNilData,
		},
		{
			name: "same upd time",
			prepare: func() []*models.Data {
				return []*models.Data{models.TestData(t)}
			},
			arg: func() *models.Data {
				data := models.TestData(t)

				return data
			},
			exp: store.ErrOldData,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, teardown := sqlstore.TestStore(t)
			defer teardown(DATA_TABLE)

			if tc.prepare != nil {
				for _, data := range tc.prepare() {
					err := s.Data().Save(data)
					assert.NoError(t, err)
				}
			}

			err := s.Data().Save(tc.arg())
			assert.Equal(t, tc.exp, err)
		})
	}
}

func TestDataRepo_Latest(t *testing.T) {
	s, teardown := sqlstore.TestStore(t)
	defer teardown(DATA_TABLE)

	testData := models.TestData(t)

	t.Run("empty", func(t *testing.T) {
		data, err := s.Data().Latest()
		assert.Nil(t, data)
		assert.NoError(t, err)
	})

	t.Run("not_empty", func(t *testing.T) {
		err := s.Data().Save(testData)
		assert.NoError(t, err)

		data, err := s.Data().Latest()
		assert.NotNil(t, data)
		assert.NoError(t, err)
		assert.ObjectsAreEqualValues(testData, data)
	})
}

func TestDataRepo_Len(t *testing.T) {
	s, teardown := sqlstore.TestStore(t)
	defer teardown(DATA_TABLE)

	t.Run("zero", func(t *testing.T) {
		len, err := s.Data().Len()
		assert.NoError(t, err)
		assert.Equal(t, 0, len)
	})

	t.Run("one", func(t *testing.T) {
		err := s.Data().Save(models.TestData(t))
		assert.NoError(t, err)

		len, err := s.Data().Len()
		assert.NoError(t, err)
		assert.Equal(t, 1, len)
	})
}

func TestDataRepo_Cleanup(t *testing.T) {
	testCases := []struct {
		name        string
		input       func() []*models.Data
		cleanupTime time.Time
		expectedLen int
	}{
		{
			name: "no data",
			input: func() []*models.Data {
				return []*models.Data{}
			},
			cleanupTime: time.Now(),
			expectedLen: 0,
		}, {
			name: "with one element.delete",
			input: func() []*models.Data {
				testData := models.TestData(t)
				testData.Upd -= 1
				return []*models.Data{testData}
			},
			cleanupTime: time.Now(),
			expectedLen: 0,
		}, {
			name: "with one element.keep",
			input: func() []*models.Data {
				testData := models.TestData(t)
				return []*models.Data{testData}
			},
			cleanupTime: time.Now().Add(-time.Second),
			expectedLen: 1,
		}, {
			name: "with two elements",
			input: func() []*models.Data {
				testData1 := models.TestData(t)
				testData1.Upd -= 60 * 60
				testData2 := models.TestData(t)
				testData2.Upd -= 1

				return []*models.Data{testData1, testData2}
			},
			cleanupTime: time.Now(),
			expectedLen: 0,
		}, {
			name: "with many elements.keep last",
			input: func() []*models.Data {
				testData1 := models.TestData(t)
				testData1.Upd -= 2 * 60 * 60
				testData2 := models.TestData(t)
				testData2.Upd -= 60 * 60
				testData3 := models.TestData(t)
				testData3.Upd -= 1

				return []*models.Data{testData1, testData2, testData3}
			},
			cleanupTime: time.Now().Add(-2 * time.Second),
			expectedLen: 1,
		}, {
			name: "with many elements.keep last 2",
			input: func() []*models.Data {
				testData1 := models.TestData(t)
				testData1.Upd -= 2 * 60 * 60
				testData2 := models.TestData(t)
				testData2.Upd -= 60 * 60
				testData3 := models.TestData(t)
				testData3.Upd -= 4
				testData4 := models.TestData(t)
				testData4.Upd -= 1
				testData5 := models.TestData(t)

				return []*models.Data{testData1, testData2, testData3, testData4, testData5}
			},
			cleanupTime: time.Now().Add(-2 * time.Second),
			expectedLen: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, teardown := sqlstore.TestStore(t)
			defer teardown(DATA_TABLE)

			args := tc.input()
			numArgs := len(args)
			for _, data := range args {
				err := s.Data().Save(data)
				assert.NoError(t, err)
			}

			affected, err := s.Data().Cleanup(tc.cleanupTime)
			assert.NoError(t, err)
			assert.Equal(t, int64(numArgs-tc.expectedLen), affected)

			len, err := s.Data().Len()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedLen, len)
		})
	}
}

func BenchmarkDataRepo(b *testing.B) {
	b.ReportAllocs()

	store, teardown := sqlstore.TestStore(b)
	defer teardown(DATA_TABLE)

	repo := store.Data()

	users := make([]models.Data, 0, b.N)
	for i := 0; i < b.N; i++ {
		users = append(users, *models.TestData(b))
	}

	zerolog.SetGlobalLevel(zerolog.Disabled)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo.Save(&users[i])
		repo.Latest()
	}
}
