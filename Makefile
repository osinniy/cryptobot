test:
	go test -v -covermode=count -coverprofile=coverage.out -timeout 30s ./...

sec:
	gosec -fmt=sarif -out=gosec.sarif -no-fail cmd/... internal/...
	gosec -no-fail cmd/... internal/...

build:
ifdef version
	docker build --build-arg ip=$(ip) -t osinniy/cryptobot:$(version) .
else
	docker build --build-arg ip=$(ip) -t osinniy/cryptobot:dev .
endif

run:
ifdef port
ifdef version
	docker run -d -p 2121:2121 -p $(port):$(port) -v cbot-files:/app/files --name cryptobot osinniy/cryptobot:$(version)
else
	docker run -d -p 2121:2121 -p $(port):$(port) -v cbot-files:/app/files --name cryptobot osinniy/cryptobot:dev
endif
else
ifdef version
	docker run -d -p 2121:2121 -v cbot-files:/app/files --name cryptobot osinniy/cryptobot:$(version)
else
	docker run -d -p 2121:2121 -v cbot-files:/app/files --name cryptobot osinniy/cryptobot:dev
endif
endif

start:
	docker start cryptobot

stop:
	docker stop cryptobot

rm:
	docker rm cryptobot

rmi:
ifdef version
	docker rmi osinniy/cryptobot:$(version)
else
	docker rmi osinniy/cryptobot:dev
endif

logs:
	docker logs -f cryptobot

push: test
ifdef version
	docker push osinniy/cryptobot:$(version)
else
	docker push osinniy/cryptobot:dev
endif

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
