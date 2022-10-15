package logs

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var slowReqThreshold uint

// WrappedLogger is just zerolog.Logger with additional helpful methods.
// type WrappedLogger zerolog.Logger

// Logs response status code, host, method, url and time elapsed.
// Warns if time elapsed is greater than slowReqThreshold.
// To change threshold use [SetSlowReqThreshold].
// Timestamp must be set right before sending request.
func Response(res *http.Response, timestamp time.Time) {
	elapsed := time.Since(timestamp)

	msg := "http request"
	logE := log.Debug()

	if elapsed.Milliseconds() > int64(slowReqThreshold) {
		msg = "slow http request"
		logE = log.Warn()
	}

	logE.Int("status", res.StatusCode).
		Str("host", res.Request.URL.Host).
		Str("method", res.Request.Method).
		Str("url", res.Request.URL.String()).
		Dur("elapsed_ms", elapsed).
		Msg(msg)
}

func SetSlowReqThreshold(threshold uint) {
	slowReqThreshold = threshold
}

func ResponseBody(data []byte) {
	log.Trace().Str("data", string(data)).Caller(1).Msg("response body")
}

// Logs query to database, elapsed time, error if occurred.
// queryName must be in format "repo.method".
func DbQuery(logger *zerolog.Logger, queryName string, err error, startTimestamp time.Time, args ...map[string]any) {
	endTime := time.Since(startTimestamp)
	logE := logger.Debug()
	if err != nil {
		logE = logger.Error().Err(err)
	}
	if len(args) > 0 {
		logE.Fields(args[0])
	}
	logE.Str("query", queryName).Int64("elapsed_µs", endTime.Microseconds()).Msg("db request")
}

// Logs memory store query, error if occurred.
// queryName must be in format "repo.method".
func MemQuery(logger *zerolog.Logger, queryName string, err error, args ...map[string]any) {
	logE := logger.Debug()
	if err != nil {
		logE = logger.Error().Err(err)
	}
	if len(args) > 0 {
		logE.Fields(args[0])
	}
	logE.Str("query", queryName).Msg("mem request")
}
