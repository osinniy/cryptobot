package app

import (
	"os"
	"os/signal"
	"osinniy/cryptobot/internal/api"
	"osinniy/cryptobot/internal/bot"
	"osinniy/cryptobot/internal/cmd"
	"osinniy/cryptobot/internal/config"
	"osinniy/cryptobot/internal/logs"
	"osinniy/cryptobot/internal/service"
	"osinniy/cryptobot/internal/store"
	"osinniy/cryptobot/internal/store/memstore"
	"osinniy/cryptobot/internal/store/sqlstore"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Run(flags cmd.Flags) {
	logs.Setup()

	// config
	cfg := config.Configure(flags.ConfigPath, true)
	if cfg == nil {
		return
	}

	if flags.Debug && zerolog.GlobalLevel() > zerolog.DebugLevel {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// store
	var store store.Store
	if flags.UseMemStore {
		store = memstore.New()
		log.Info().Msg("using in-memory store")
	} else {
		sqlstore := sqlstore.Init(cfg.Database.Path)
		store = sqlstore
		// also we can use store with caching:
		// memstore := memstore.New()
		// store = mergedstore.Merge(memstore, sqlitestore)
	}

	// service
	service := service.New(
		store.Data(),
		time.Duration(cfg.Service.RefreshInterval*1000*1000*1000), // seconds to ns
		time.Duration(cfg.Service.KeepAlive*60*60*1000*1000*1000), // hours to ns
	)
	go service.Start()

	// api
	port, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		log.Warn().Msg("failed to parse API_PORT environment variable, using default port")
	}
	server := api.NewServer(store, port)
	go server.Serve()

	// bot
	bot := bot.Setup(store, cfg.Secrets.BotToken)
	go bot.Run()

	// shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	sig := <-exit

	log.Info().Msg(sig.String())

	service.Stop()
	bot.Stop()
	err = store.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to close store")
	}
}
