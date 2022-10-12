package config

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/rs/zerolog/log"
)

type Webhook struct {
	Enabled bool
	Ip      string
	Port    int
	Secret  string
	PubKey  string
	PrivKey string
}

type Config struct {
	Secrets struct {
		BotToken  string
		CMCApiKey string
	}

	Webhook Webhook

	Database struct {
		Path string
	}

	Service struct {
		RefreshInterval uint
		KeepAlive       uint
	}

	Logs struct {
		Level            string
		SlowReqThreshold uint
		Path             string
	}
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
	validate(conf)

	return conf
}

func checkDefaults(conf *Config) {
	if conf.Webhook.Port == 0 {
		conf.Webhook.Port = DEFAULT_WEBHOOK_PORT
	}
	if conf.Webhook.PubKey == "" {
		conf.Webhook.PubKey = DEFAULT_WEBHOOK_PUBLIC_KEY
	}
	if conf.Webhook.PrivKey == "" {
		conf.Webhook.PrivKey = DEFAULT_WEBHOOK_PRIVATE_KEY
	}
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

func validate(conf *Config) {
	switch conf.Webhook.Port {
	case 0, 443, 80, 88, 8443:
	default:
		log.Error().Int("port", conf.Webhook.Port).Msg("invalid webhook port in config")
	}

	for i, char := range conf.Webhook.Secret {
		if (char < 'a' || char > 'z') &&
			(char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') &&
			char != '_' &&
			char != '-' {
			i++
			log.Error().Str("char", string(char)).Int("position", i).
				Msg("configured webhook secret contains invalid character")
			break
		}
	}
}
