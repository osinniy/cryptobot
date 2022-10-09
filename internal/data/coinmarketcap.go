package data

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"osinniy/cryptobot/internal/config"
	"osinniy/cryptobot/internal/logs"

	"github.com/rs/zerolog/log"
)

const (
	CMC_SITE_URL = "https://coinmarketcap.com/"
	CMC_API_URL  = "https://pro-api.coinmarketcap.com/"
)

// Returns latest global market stats from coinmarketcap.com
//
// # Call it no more than 5 minutes
//
// Endpoint: /v1/global-metrics/quotes/latest
func LatestMarketStats() (result *CMCMetricsResponse, err error) {
	req, err := http.NewRequest("GET", CMC_API_URL+"v1/global-metrics/quotes/latest", nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to create request")
		return
	}

	req.Header.Add("X-CMC_PRO_API_KEY", config.Global.Secrets.CMCApiKey)
	timestamp := time.Now()
	res, err := httpClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("failed to send request")
		return
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return
	}
	logs.Response(res, timestamp)
	logs.ResponseBody(data)

	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode CMC metrics response")
		return
	}

	if result.Status.ErrorCode != 0 {
		err = errors.New(result.Status.ErrorMessage)
		log.Error().
			Err(err).
			Int("code", result.Status.ErrorCode).
			Msg("failed to obtain market stats: CMC returned non-successfull response")
		return
	}

	return
}
