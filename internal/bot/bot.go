package bot

import (
	"io"
	"osinniy/cryptobot/internal/config"
	"osinniy/cryptobot/internal/store"
	"strconv"
	"strings"
	"time"

	stdLog "log"

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

// Initiates bot, setups handlers and buttons so it can be started.
// Uses long polling or webhook id set.
func Setup(store store.Store, token string, webhook config.Webhook) (b *Bot) {
	b = &Bot{
		store:  store,
		logger: log.With().Str("module", MODULE).Logger(),
	}

	updates := []string{"message", "callback_query"}
	settings := tg.Settings{
		Token:   token,
		OnError: b.onError,
		Poller: &tg.LongPoller{
			AllowedUpdates: updates,
		},
		Verbose: true,
	}

	// We use this call because without it bot verbose mode
	// will produce a lot of garbage in logs. We need to use
	// verbose to be able to handle errors while starting
	// polling and setting webhook. But it disables standard
	// logs globally so use it with caution.
	stdLog.SetOutput(io.Discard)

	timestamp := time.Now()
	var err error
	b.tbot, err = tg.NewBot(settings)
	if err != nil {
		b.logger.Fatal().Err(err).Str("token", token).Msg("failed to setup bot")
	}

	if webhook.Enabled {
		b.tbot.Poller = &tg.Webhook{
			Listen:      ":" + strconv.Itoa(webhook.Port),
			SecretToken: webhook.Secret,
			TLS: &tg.WebhookTLS{
				Cert: webhook.PubKey,
				Key:  webhook.PrivKey,
			},
			Endpoint: &tg.WebhookEndpoint{
				PublicURL: "https://" + webhook.Ip + ":" + strconv.Itoa(webhook.Port),
				Cert:      webhook.PubKey,
			},
			AllowedUpdates: updates,
		}
	} else {
		err := b.tbot.RemoveWebhook()
		if err != nil {
			b.logger.Error().Err(err).Msg("failed to remove webhook, polling won't work")
		}
		b.logger.Warn().Msg("webhook is not configured or disabled, using long polling")
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
	defer func() {
		p := recover()
		if p != nil {
			b.logger.Panic().Msgf("Unexpected bot shutdown: %v", p)
		}
	}()

	b.running = true
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
	if strings.Contains(err.Error(), "Post") && strings.Contains(err.Error(), "context canceled") {
		return // filter getUpdates request canceling after bot shutdown in verbose mode
	}
	if err == tg.ErrSameMessageContent {
		return
	}

	logE := b.logger.Error().Err(err).Stack()
	if c != nil {
		logE.Int64("chat_id", c.Chat().ID)
	}
	logE.Msg("bot error occurred")
}

func (b Bot) logUpdate(c tg.Context) {
	b.logger.Debug().
		Int64("chat_id", c.Chat().ID).
		Int("upd_id", c.Update().ID).
		Str("text", c.Text()).
		Msg("incoming update")
}
