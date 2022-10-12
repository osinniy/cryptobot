package bot

import (
	"osinniy/cryptobot/internal/store"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tg "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

const MODULE = "bot"

type Bot struct {
	running bool
	tbot    *tg.Bot
	store   store.Store
	logger  zerolog.Logger
}

// Initiates bot, setups handlers and buttons so it can be started
func Setup(store store.Store, token string) (b *Bot) {
	b = &Bot{
		store:  store,
		logger: log.With().Str("module", MODULE).Logger(),
	}

	settings := tg.Settings{
		Token:   token,
		OnError: b.onError,
		Poller: &tg.LongPoller{
			AllowedUpdates: []string{"message", "callback_query"},
		},
	}
	if zerolog.GlobalLevel() == zerolog.TraceLevel {
		settings.Verbose = true
	}

	timestamp := time.Now()
	var err error
	b.tbot, err = tg.NewBot(settings)
	if err != nil {
		b.logger.Fatal().Err(err).Str("token", token).Msg("failed to setup bot")
	}

	b.logger.Info().
		Int64("bot_id", b.tbot.Me.ID).
		Str("bot_username", b.tbot.Me.Username).
		Bool("groups", b.tbot.Me.CanJoinGroups).
		Bool("inline", b.tbot.Me.SupportsInline).
		Dur("elapsed_ms", time.Since(timestamp)).
		Msg("bot setup completed")

	// Log all incoming updates
	b.tbot.Use(func(next tg.HandlerFunc) tg.HandlerFunc {
		return func(c tg.Context) error {
			b.logUpdate(c)
			return next(c)
		}
	})

	// Handle panics. Without arguments it will call [Bot.onError] with nil context
	b.tbot.Use(middleware.Recover())

	b.initHandlers()
	b.initButtons()

	return
}

// Start consuming updates
func (b *Bot) Run() {
	b.running = true
	b.logger.Info().Msg("waiting for updates")
	b.tbot.Start()
}

// Gracefull shutdown
func (b *Bot) Stop() {
	if !b.running {
		return
	}

	b.running = false
	b.tbot.Stop()
}

func (b *Bot) onError(err error, c tg.Context) {
	if err == tg.ErrSameMessageContent {
		return
	}

	logE := b.logger.Error().Err(err).Stack()
	if c != nil {
		logE.Int64("chat_id", c.Chat().ID)
	}
	logE.Msg("handling incoming update error occured")
}

func (b Bot) logUpdate(c tg.Context) {
	b.logger.Debug().
		Int64("chat_id", c.Chat().ID).
		Int("upd_id", c.Update().ID).
		Str("text", c.Text()).
		Msg("incoming update")
}
