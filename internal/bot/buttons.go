package bot

import tg "gopkg.in/telebot.v3"

var (
	btnsStart   = &tg.ReplyMarkup{}
	btnsStats   = &tg.ReplyMarkup{}
	btnsFutures = &tg.ReplyMarkup{}

	btnStats          = btnsStart.Data("Stats", "stats")
	btnFutures        = btnsStart.Data("Futures", "futures")
	btnRefreshStats   = btnsStats.Data("🔄", "refresh_stats")
	btnRefreshFutures = btnsFutures.Data("🔄", "refresh_futures")
)

// Call it before any usage of buttons
func (b *Bot) initButtons() {
	btnsStart.Inline(
		btnsStart.Row(btnStats, btnFutures),
	)
	btnsStats.Inline(
		btnsStats.Row(btnRefreshStats, btnFutures),
	)
	btnsFutures.Inline(
		btnsFutures.Row(btnStats, btnRefreshFutures),
	)

	b.tbot.Handle(&btnStats, b.onGlobalStats)
	b.tbot.Handle(&btnFutures, b.onFuturesStats)
	b.tbot.Handle(&btnRefreshStats, b.onGlobalStats)
	b.tbot.Handle(&btnRefreshFutures, b.onFuturesStats)
}
