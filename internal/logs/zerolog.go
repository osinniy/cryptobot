package logs

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

const (
	CONSOLE_TIME_FORMAT  = "2006-01-02 15:04:05"
	LOG_FILE_TIME_FORMAT = "2006-01-02T15:04:05"
)

var console *zerolog.ConsoleWriter

// Setups console global logger
// To change level use [SetLevel]
// To write to file additionally use [SetLogFile]
func Setup() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = LOG_FILE_TIME_FORMAT
	zerolog.LevelFieldName = "l"
	zerolog.ErrorFieldName = "err"
	zerolog.MessageFieldName = "msg"
	zerolog.DurationFieldInteger = true

	consoleWr := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = CONSOLE_TIME_FORMAT
	})
	log.Logger = log.Output(consoleWr)
	console = &consoleWr
}

// Tries to open log file and set it as global log output
func SetLogFile(file string) {
	logFile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		log.Error().Err(err).Msg("failed to open log file")
		return
	}

	if console == nil {
		log.Logger = log.Output(logFile)
	} else {
		log.Logger = log.Output(zerolog.MultiLevelWriter(console, logFile))
	}
}

// Tries to parse level and set it as global log level
func SetLevel(level string) (ok bool) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse log level")
		return
	}

	zerolog.SetGlobalLevel(lvl)
	ok = true
	return
}
