package service_test

import (
	"osinniy/cryptobot/internal/logs"
	"osinniy/cryptobot/internal/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	logs.Setup()
}

func TestService(t *testing.T) {
	t.Skip()

	srv := service.TestService(t)

	go srv.Start()
}

func TestService_IsRunning(t *testing.T) {
	srv := service.TestService(t)

	time.AfterFunc(time.Millisecond, func() {
		assert.True(t, srv.IsRunning())

		srv.Stop()
		assert.False(t, srv.IsRunning())
	})

	time.AfterFunc(10*time.Millisecond, func() {
		t.Fail()
	})
	srv.Start()
}
