# Cryptobot
![Coverage](https://img.shields.io/badge/Coverage-74.1%25-brightgreen)

![Project Schema](schema.svg)

### **Cryptobot** basically consists of three parts: Service, Store and Bot.

**[Service](internal/service/service.go)** fetches data from **[Data](internal/data)** endpoints and saves it to **[DataRepository](internal/store/repository.go)**.
**[Store](internal/store/store.go)** of these repository have third different implementations: **[SqlStore](internal/store/sqlstore)**, **[MemStore](internal/store/memstore)** and **[MergedStore](internal/store/mergedstore)** (which uses both of them).
**[Bot](internal/bot)** fetches data from store by user request received from *Telegram*.
It also have optional **[api](internal/api)** that makes accessible some metrics to *CLI*. *Prometeus* uses this cli to fetch metrics and show them in *Grafana*.
*CLI - Prometeus - Grafana* chain is not implemented yet.

## Run

To run the project follow the steps below:

### 1. Clone the project

```bash
git clone osinniy/cryptobot
cd cryptobot
```

### 2. Configure

To configure bot create `bot.yml` as shown below. Set variables to your own:

```yaml
secrets:
  # Your bot token. How to obtain your token: https://core.telegram.org/bots/features#botfather
  botToken: 0000000000:XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
  # CMC api key: https://pro.coinmarketcap.com/account
  cmcApiKey: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

It's also recommended setting up a webhook. To proceed go to [example config](configs/example/bot.yml) and fill webhook config as shown there.
Full list of variables and their default values you can also find here: [configs/example/bot.yml](configs/example/bot.yml).

### 3. Build & run

To run the project, you need to have [Docker](https://docker.com/) installed on your machine.

First, move config to `configs/release` folder:

```bash
mkdir configs/release
mv bot.yml configs/release/bot.yml
```

Then run the build. If you use webhook, pass your IP as shown below. Certificates will be generated during the build with your IP.

```bash
make ip=0.0.0.0 # Or simply make if webhook is disabled
```

After process is finished, you will have osinniy/cryptobot image. Now you can run the project.
Webhook users should also pass port they set in config so docker can expose it:

```bash
make run port=8443 # Or simply make run if webhook is disabled
```

### Manual build

If you don't want to use Docker, you can build the project manually. It requires [Go 1.19](https://golang.org/) installed.

```bash
cd build && make
```

After that you need to run migrations. We use [goose](https://github.com/pressly/goose) for them. In order to install it, use command:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

It's might be neccessary to add GOBIN to your PATH:

```bash
PATH=$PATH:$(go env GOPATH)/bin
echo $PATH >> ~/.profile
```

And then execute migrations. You must be in root project folder:

```bash
goose -dir sql sqlite3 cryptobot.sqlite up
```

And now you can run the binary:

```bash
./cryptobot
```

To see available options, use `./cryptobot --help`

## Test

For tests use:

```bash
make test
```

Note that you need to have goose installed and GOBIN added to path to be able to run tests:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest

PATH=$PATH:$(go env GOPATH)/bin
echo $PATH >> ~/.profile
```
