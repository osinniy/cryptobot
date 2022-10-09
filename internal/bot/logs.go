package bot

import (
	"github.com/rs/zerolog"
	tg "gopkg.in/telebot.v3"
)

func logNewUser(logger *zerolog.Logger, user *tg.User) {
	logger.Info().
		Int64("user_id", user.ID).
		Str("username", user.Username).
		Str("lang", user.LanguageCode).
		Msg("new user seen")
}
