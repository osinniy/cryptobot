build:
	go build -v -ldflags="-s -w" -trimpath -o cryptobot ../cmd/cryptobot/main.go

start:
	cryptobot --config ../configs/debug/bot.yml

run:
	go run ../cmd/main.go --config ../configs/debug/bot.yml

help:
	cryptobot --help

.DEFAULT_GOAL := build
