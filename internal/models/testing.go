package models

import (
	"testing"
	"time"
)

// Test user object
func TestUser(t testing.TB) *User {
	t.Helper()

	return &User{
		Id:        1,
		FirstSeen: time.Now().Unix(),
		Lang:      "en",
	}
}

// Test data object
func TestData(t testing.TB) *Data {
	t.Helper()

	return &Data{
		BTCDom:                   0,
		ETHDom:                   0,
		BTCDom24HChange:          -0.34,
		ETHDom24HChange:          0.01,
		StablecoinsCap:           149_843_547_960,
		TotalCap:                 946_438_273_998,
		TotalCap24HChange:        0.35,
		Liquidation24HNum:        64_622,
		Liquidation24HUsd:        149_107_938,
		Liquidations24HUsdChange: 39.65,
		OpenInterest:             27_793_329_789,
		OpenInterest24HChange:    0.54,
		Upd:                      time.Now().Unix(),
	}
}
