test:
	go test -v -covermode=count -coverprofile=coverage.out -timeout 30s ./...

sec:
	gosec -fmt=sarif -out=gosec.sarif -no-fail cmd/... internal/...
	gosec -no-fail cmd/... internal/...

build:
	docker build --build-arg ip=$(ip) -t osinniy/cryptobot:dev .

run:
ifdef port
	docker run -d -p 2121:2121 -p $(port):$(port) -v cbot-files:/app/files --name cryptobot osinniy/cryptobot:dev
else
	docker run -d -p 2121:2121 -v cbot-files:/app/files --name cryptobot osinniy/cryptobot:dev
endif

start:
	docker start cryptobot

stop:
	docker stop cryptobot

rm:
	docker rm cryptobot

rmi:
	docker rmi osinniy/cryptobot:dev

logs:
	docker logs -f cryptobot

push: test
	docker push osinniy/cryptobot:dev

restart:
	make stop
	make start

rebuild:
	make rm
	make rmi
	make build

rerun:
	make stop
	make rebuild
	make run

.PHONY: build

.DEFAULT_GOAL := build
