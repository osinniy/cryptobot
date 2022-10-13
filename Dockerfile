FROM golang:1.19.2-alpine3.16 AS builder

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

# Build

COPY . .

RUN go build -v -ldflags="-s -w" -trimpath -o cryptobot cmd/cryptobot/main.go

# Final image

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /src/cryptobot .
COPY --from=builder /src/configs/release/bot.yml .
COPY --from=builder /src/cryptobot.sqlite files/cryptobot.sqlite

ENV API_PORT 2121
EXPOSE $API_PORT

VOLUME [ "/app/files" ]

CMD ["/app/cryptobot"]
