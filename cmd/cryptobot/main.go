package main

import (
	"osinniy/cryptobot/internal/app"
	"osinniy/cryptobot/internal/cmd"
)

func main() {
	flags := cmd.ParseFlags()

	app.Run(flags)
}
