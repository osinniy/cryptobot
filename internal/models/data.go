package models

import (
	"osinniy/cryptobot/internal/data"
	"time"
)

type Data struct {
	BTCDom                   float64
	ETHDom                   float64
	BTCDom24HChange          float64
	ETHDom24HChange          float64
	StablecoinsCap           float64
	TotalCap                 float64
	TotalCap24HChange        float64
	Liquidation24HNum        uint32
	Liquidation24HUsd        float64
	Liquidations24HUsdChange float32
	OpenInterest             uint64
	OpenInterest24HChange    float32
	Upd                      int64
}

// Constructs a new Data from cmc and coinglass responses.
func BuildData(cmc data.CMCMetricsResponse, coinglass data.CoinglassResponse) *Data {
	return &Data{
		BTCDom:                   cmc.Data.BTCDom,
		ETHDom:                   cmc.Data.ETHDom,
		BTCDom24HChange:          cmc.Data.BTCDom24HChange,
		ETHDom24HChange:          cmc.Data.ETHDom24HChange,
		StablecoinsCap:           cmc.Data.StablecoinsCap,
		TotalCap:                 cmc.Data.Quote.USD.TotalCap,
		TotalCap24HChange:        cmc.Data.Quote.USD.TotalCap24HChange,
		Liquidation24HNum:        coinglass.Data.Liquidation24HNum,
		Liquidation24HUsd:        coinglass.Data.Liquidation24HUsd,
		Liquidations24HUsdChange: coinglass.Data.Liquidations24HUsdChange,
		OpenInterest:             coinglass.Data.OpenInterest,
		OpenInterest24HChange:    coinglass.Data.OpenInterest24HChange,
		Upd:                      time.Now().Unix(),
	}
}
