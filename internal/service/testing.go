package service

import (
	"osinniy/cryptobot/internal/store/memstore"
	"testing"
	"time"
)

// Test service object
func TestService(t *testing.T) *Service {
	t.Helper()

	return New(memstore.New().Data(), 5*time.Second, time.Minute)
}
