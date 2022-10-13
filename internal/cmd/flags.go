package cmd

import "flag"

type Flags struct {
	Debug       bool
	ConfigPath  string
	UseMemStore bool
	DbPath      string
	CertPath    string
	KeyPath     string
}

func ParseFlags() (flags Flags) {
	flag.StringVar(&flags.ConfigPath, "config", "bot.yml", "path to config file")
	flag.BoolVar(&flags.UseMemStore, "mem", false, "use in-memory temporary store")
	flag.BoolVar(&flags.Debug, "debug", false, "enable debug logs. It will override config log level")
	flag.StringVar(&flags.DbPath, "db", "", "path to database file. Overrides config")
	flag.StringVar(&flags.CertPath, "cert", "", "path to certificate file. Overrides config")
	flag.StringVar(&flags.KeyPath, "key", "", "path to key file. Overrides config")

	flag.Parse()
	return
}
