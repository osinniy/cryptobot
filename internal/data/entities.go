package data

// Standardized coinmarketcap status object
type CMCStatus struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      int    `json:"elapsed"`
	CreditCount  int    `json:"credit_count"`
}

// Latest global cryptocurrency market metrics
//
// Req: https://pro-api.coinmarketcap.com/v1/global-metrics/quotes/latest
//
// API Ref: https://coinmarketcap.com/api/documentation/v1/#operation/getV1GlobalmetricsQuotesLatest
type CMCMetricsResponse struct {
	Status CMCStatus `json:"status"`
	Data   struct {
		ActiveCryptos           int     `json:"active_cryptocurrencies"`
		TotalCryptos            int     `json:"total_cryptocurrencies"`
		ActiveMarketPairs       int     `json:"active_market_pairs"`
		ActiveExchanges         int     `json:"active_exchanges"`
		TotalExchanges          int     `json:"total_exchanges"`
		BTCDom                  float64 `json:"btc_dominance"`
		ETHDom                  float64 `json:"eth_dominance"`
		BTCDom24HChange         float64 `json:"btc_dominance_24h_percentage_change"`
		ETHDom24HChange         float64 `json:"eth_dominance_24h_percentage_change"`
		DefiVol24H              float64 `json:"defi_volume_24h"`
		DefiVol24HReport        float64 `json:"defi_volume_24h_reported"`
		DefiVol24HChange        float64 `json:"defi_24h_percentage_change"`
		DefiCap                 float64 `json:"defi_market_cap"`
		StablecoinsVol24H       float64 `json:"stablecoin_volume_24h"`
		StablecoinsVol24HReport float64 `json:"stablecoin_volume_24h_reported"`
		StablecoinsVol24HChange float64 `json:"stablecoin_24h_percentage_change"`
		StablecoinsCap          float64 `json:"stablecoin_market_cap"`
		DerivativesVol24H       float64 `json:"derivatives_volume_24h"`
		DerivativesVol24HReport float64 `json:"derivatives_volume_24h_reported"`
		DerivativesVol24HChange float64 `json:"derivatives_24h_percentage_change"`
		LastUpdated             string  `json:"last_updated"`

		Quote struct {
			USD struct {
				TotalCap             float64 `json:"total_market_cap"`
				TotalCapYesterday    float64 `json:"total_market_cap_yesterday"`
				TotalCap24HChange    float64 `json:"total_market_cap_yesterday_percentage_change"`
				TotalVol24H          float64 `json:"total_volume_24h"`
				TotalVol24HReport    float64 `json:"total_volume_24h_reported"`
				TotalVolYesterday    float64 `json:"total_volume_24h_yesterday"`
				TotalVol24HChange    float64 `json:"total_volume_24h_yesterday_percentage_change"`
				AltcoinsVol24H       float64 `json:"altcoin_volume_24h"`
				AltcoinsVol24HReport float64 `json:"altcoin_volume_24h_reported"`
				AltcoinsCap          float64 `json:"altcoin_market_cap"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}

// type CMCCategoryData struct {
// 	Id              string  `json:"id"`
// 	Name            string  `json:"name"`
// 	Title           string  `json:"title"`
// 	Description     string  `json:"description"`
// 	NumTokens       int     `json:"num_tokens"`
// 	AvgPriceChange  float64 `json:"avg_price_change"`
// 	MarketCap       float64 `json:"market_cap"`
// 	MarketCapChange float64 `json:"market_cap_change"`
// 	Volume          float64 `json:"volume"`
// 	VolumeChange    float64 `json:"volume_change"`
// 	LastUpdated     string  `json:"last_updated"`
// }

// type CMCResponse struct {
// 	Data   CMCCategoryData `json:"data"`
// 	Status CMCStatus       `json:"status"`
// }

type CoinglassData struct {
	Liquidation24HNum        uint32  `json:"liquidationH24Num"`
	Liquidation24HUsd        float64 `json:"liquidationH24VolUsd"`
	Liquidations24HUsdChange float32 `json:"lqH24Chain"`
	LongRate                 float32 `json:"longRate"`
	ShortRate                float32 `json:"shortRate"`
	Volume24HUsd             uint64  `json:"volUsd"`
	Volume24HChange          float32 `json:"volH24Chain"`
	OpenInterest             uint64  `json:"openInterest"`
	OpenInterest24HChange    float32 `json:"oiH24Chain"`
}

type CoinglassResponse struct {
	Code    string        `json:"code"`
	Msg     string        `json:"msg"`
	Data    CoinglassData `json:"data"`
	Success bool          `json:"success"`
}
