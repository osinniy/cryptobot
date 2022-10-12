package bot_test

import (
	"osinniy/cryptobot/internal/bot"
	"osinniy/cryptobot/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FmtGlobalStats(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		data func(*testing.T) *models.Data
		want string
	}{
		{
			name: "empty data",
			data: func(t *testing.T) *models.Data {
				return &models.Data{}
			},
			want: `*Global Market Stats:*

Total Market Cap: *$0 (+0,00%)*
Stablecoins Cap: *$0*
BTC Dominance: *0,0% (+0,00%)*
ETH Dominance: *0,0% (+0,00%)*

Last update: _ 1 Jan 00:00:00 UTC_
`,
		},
		{
			name: "normal data",
			data: func(t *testing.T) *models.Data {
				data := models.TestData(t)
				data.BTCDom = 39.8
				data.ETHDom = 17.1
				data.Upd = 1665544708
				return data
			},
			want: `*Global Market Stats:*

📈Total Market Cap: *$946 438 273 998 (+0,35%)*
Stablecoins Cap: *$149 843 547 960*
📉BTC Dominance: *39,8% (-0,34%)*
📈ETH Dominance: *17,1% (+0,01%)*

Last update: _12 Oct 03:18:28 UTC_
`, // has nbsp in numbers
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := bot.FmtGlobalStats(*tc.data(t))
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_FmtFuturesStats(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		data func(*testing.T) *models.Data
		want string
	}{
		{
			name: "empty data",
			data: func(t *testing.T) *models.Data {
				return &models.Data{}
			},
			want: `*Futures Stats:*

Open Interest: *$0 (+0,00%)*
Liquidations 24H: *$0 (+0,00%)*
Number of liquidations: *0*

Last update: _ 1 Jan 00:00:00 UTC_
`,
		},
		{
			name: "normal data",
			data: func(t *testing.T) *models.Data {
				data := models.TestData(t)
				data.Upd = 1665544708
				return data
			},
			want: `*Futures Stats:*

📈Open Interest: *$27 793 329 789 (+0,54%)*
📈Liquidations 24H: *$149 107 938 (+39,65%)*
Number of liquidations: *64 622*

Last update: _12 Oct 03:18:28 UTC_
`, // has nbsp in numbers
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := bot.FmtFuturesStats(*tc.data(t))
			assert.Equal(t, tc.want, got)
		})
	}
}
