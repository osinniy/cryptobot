FROM golang:1.20.0-alpine3.16 AS builder

WORKDIR /src

# Install gcc

RUN apk update && apk upgrade
RUN apk add --update gcc musl-dev

# Download dependencies

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations

COPY sql ./

RUN goose -v sqlite3 cryptobot.sqlite up

# Generate certificate

ARG ip

RUN apk add openssl
RUN openssl req -newkey rsa:2048 -sha256 -nodes -keyout key.pem -x509 -days 365 -out cert.pem -subj "/C=/ST=/L=/O=/CN=$ip"

# Build

COPY . .

RUN go build -v -ldflags="-s -w" -trimpath -o cryptobot cmd/cryptobot/main.go

# Final image

FROM alpine:3.17

WORKDIR /app

COPY --from=builder /src/cryptobot .
COPY --from=builder /src/configs/release/bot.yml .
COPY --from=builder /src/cryptobot.sqlite files/cryptobot.sqlite
COPY --from=builder /src/cert.pem files/cert.pem
COPY --from=builder /src/key.pem files/key.pem

# Install certificates

RUN apk update --no-cache && apk upgrade --no-cache
RUN apk add --no-cache ca-certificates

ENV API_PORT 2121
EXPOSE $API_PORT

VOLUME [ "/app/files" ]

CMD ["/app/cryptobot", "-db", "files/cryptobot.sqlite", "-cert", "files/cert.pem", "-key", "files/key.pem"]
