package bot

import (
	"osinniy/cryptobot/internal/models"
	"time"

	"github.com/rs/zerolog"
	tg "gopkg.in/telebot.v3"
)

func (b *Bot) initHandlers() {
	b.tbot.Handle("/start", b.onStart)
	b.tbot.Handle("/help", b.onHelp)
}

// Command: /start
func (b *Bot) onStart(c tg.Context) error {
	user := c.Sender()

	b.logger.Debug().Func(func(e *zerolog.Event) {
		// automatically logs exists state
		b.store.Users().IsExists(user.ID)
	})

	err := b.store.Users().Add(&models.User{
		Id:        user.ID,
		FirstSeen: time.Now().Unix(),
		Lang:      user.LanguageCode,
	})
	if err != nil {
		return nil
	}

	logNewUser(&b.logger, user)

	return c.Send("Hello, "+user.FirstName+"!. Welcome to CryptoBot", btnsStart)
}

// Command: /help
func (b *Bot) onHelp(c tg.Context) error {
	return c.Send("Available commands: /start, /help")
}

// Button: stats
func (b *Bot) onGlobalStats(c tg.Context) error {
	data, err := b.getData(c)
	if err != nil || data == nil {
		return err
	}
	_, err = b.tbot.Edit(c.Message(), fmtGlobalStats(*data), tg.ModeMarkdown, btnsStats)
	return err
}

// Button: futures
func (b *Bot) onFuturesStats(c tg.Context) error {
	data, err := b.getData(c)
	if err != nil || data == nil {
		return err
	}

	_, err = b.tbot.Edit(c.Message(), fmtFuturesStats(*data), tg.ModeMarkdown, btnsFutures)
	return err
}

func (b *Bot) getData(c tg.Context) (data *models.Data, err error) {
	data, err = b.store.Data().Latest()
	if err != nil || data == nil {
		return nil, c.Send("No data yet. Retry after a while.")
	}
	return data, nil
}
