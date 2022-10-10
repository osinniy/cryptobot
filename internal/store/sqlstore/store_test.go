package sqlstore_test

import (
	"fmt"
	"os"
	"os/exec"
	"osinniy/cryptobot/internal/logs"
	"osinniy/cryptobot/internal/store/sqlstore"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ROOT_PATH = sqlstore.ROOT_PATH
	DB_PATH   = sqlstore.DB_PATH
)

func TestMain(m *testing.M) {
	logs.Setup()

	prepare()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func prepare() {
	if _, err := os.Create(DB_PATH); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := exec.Command("goose", "-dir", ROOT_PATH+"sql", "sqlite3", DB_PATH, "up")
	out, err := cmd.CombinedOutput()
	fmt.Println()
	fmt.Println(string(out))

	if err != nil {
		fmt.Println(err)
	}
}

func teardown() {
	if err := os.Remove(DB_PATH); err != nil {
		fmt.Println(err)
	}
}

func TestInvalidStore(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skip test on windows")
	}

	invalidPath := "/invalid_db"
	s, db, err := sqlstore.TestDb(t, invalidPath)
	assert.Nil(t, db)
	assert.Equal(t, err.Error(), "unable to open database file: no such file or directory")
	assert.Equal(t, invalidPath, s.Path)
}
