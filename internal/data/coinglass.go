package data

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"osinniy/cryptobot/internal/logs"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	COINGLASS_URL     = "https://coinglass.com/"
	COINGLASS_API_URL = "https://fapi.coinglass.com/"
	USER_AGENT        = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"
)

// Returns coinglass market stats.
// Them includes futures stats, liquidations info and open interest.
//
// Uses coinglass.com undocumented api so use it carefully
func CoinglassStats() (result *CoinglassResponse, err error) {
	req, err := http.NewRequest("GET", COINGLASS_API_URL+"api/futures/home/statistics", nil)
	if err != nil {
		log.Error().Msg("failed to create request")
		return
	}

	req.Header.Add("User-Agent", USER_AGENT)
	timestamp := time.Now()
	res, err := httpClient.Do(req)
	if err != nil {
		log.Error().Msg("failed to send request")
		return
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().Msg("failed to read response body")
		return
	}
	logs.Response(res, timestamp)
	logs.ResponseBody(data)

	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Error().Msg("failed to decode coinglass response")
		return
	}
	if !result.Success {
		err = errors.New(result.Msg)
		log.Error().
			Err(err).
			Str("code", result.Code).
			Msg("failed to obtain market liquidations: coinglass returned non-successful response")
		return
	}

	return
}

// TODO: add endpoint
// https://fapi.coinglass.com/api/futures/liquidation/order?side=&exName=&symbol=&pageSize=100&pageNum=1&volUsd=1000000
