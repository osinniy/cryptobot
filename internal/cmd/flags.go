package cmd

import "flag"

type Flags struct {
	Debug       bool
	ConfigPath  string
	UseMemStore bool
}

func ParseFlags() (flags Flags) {
	flag.BoolVar(&flags.Debug, "debug", false, "enable debug logs. It will override config log level")
	flag.StringVar(&flags.ConfigPath, "config", "bot.yml", "path to config file")
	flag.BoolVar(&flags.UseMemStore, "mem", false, "use in-memory temporary store")

	flag.Parse()
	return
}
