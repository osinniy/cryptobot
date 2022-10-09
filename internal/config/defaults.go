package config

const (
	DEFAULT_DATABASE_PATH = "cryptobot.sqlite"

	DEFAULT_SERVICE_REFRESH_INTERVAL = 600 // 10 minutes
	DEFAULT_SERVICE_KEEP_ALIVE_TIME  = 168 // 1 week

	DEFAULT_LOG_LEVEL              = "info"
	DEFAULT_LOG_SLOW_REQ_THRESHOLD = 1500 // milliseconds
)
