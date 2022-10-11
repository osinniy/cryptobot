package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Secrets struct {
		BotToken  string `mapstructure:"botToken"`
		CMCApiKey string `mapstructure:"cmcApiKey"`
	} `mapsstructure:"secrets"`

	Database struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"database"`

	Service struct {
		RefreshInterval uint `mapstructure:"refreshInterval"`
		KeepAlive       uint `mapstructure:"keepAlive"`
	} `mapstructure:"service"`

	Logs struct {
		Level            string `mapstructure:"level"`
		SlowReqThreshold uint   `mapstructure:"slowReqThreshold"`
		Path             string `mapstructure:"path"`
	} `mapstructure:"logs"`
}

// Reads config from configFile. If global is true, then saves config to config.Global
func Configure(configFile string) *Config {
	cfg := config.NewWith(configFile, func(c *config.Config) {
		c.WithOptions(config.ParseEnv).AddDriver(yamlv3.Driver)
	})

	// Load config from file
	err := cfg.LoadFiles(configFile)
	if err != nil {
		log.Error().Err(err).Str("config", configFile).Msg("failed to load config")
		return nil
	}

	// Bind config to struct
	conf := &Config{}
	err = cfg.Decode(conf)
	if err != nil {
		log.Error().Err(err).Str("config", configFile).Msg("failed to parse config")
		return nil
	}

	checkDefaults(conf)

	return conf
}

func checkDefaults(conf *Config) {
	if conf.Database.Path == "" {
		conf.Database.Path = DEFAULT_DATABASE_PATH
	}
	if conf.Service.RefreshInterval == 0 {
		conf.Service.RefreshInterval = DEFAULT_SERVICE_REFRESH_INTERVAL
	}
	if conf.Service.KeepAlive == 0 {
		conf.Service.KeepAlive = DEFAULT_SERVICE_KEEP_ALIVE_TIME
	}
	if conf.Logs.SlowReqThreshold == 0 {
		conf.Logs.SlowReqThreshold = DEFAULT_LOG_SLOW_REQ_THRESHOLD
	}
	if conf.Logs.Level == "" {
		conf.Logs.Level = DEFAULT_LOG_LEVEL
	}
}
