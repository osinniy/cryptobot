package bot

import (
	"osinniy/cryptobot/internal/models"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	TIME_FORMAT = "_2 Jan 15:04:05 UTC"

	EMJ_GRAPHUP   = "📈"
	EMJ_GRAPHDOWN = "📉"
)

// Output must be send with tg.ModeMarkdown
func fmtGlobalStats(data models.Data) string {
	var sb strings.Builder
	sb.Grow(256)
	msg := message.NewPrinter(language.Ukrainian)

	// Header
	sb.WriteString("*Global Market Stats:*\n\n")

	// Data
	wrUpDownEmj(&sb, float64(data.TotalCap24HChange))
	sb.WriteString(msg.Sprintf("Total Market Cap: *$%.f (%+.2f%%)*\n", data.TotalCap, data.TotalCap24HChange))

	sb.WriteString(msg.Sprintf("Stablecoins Cap: *$%.f*\n", data.StablecoinsCap))

	wrUpDownEmj(&sb, float64(data.BTCDom24HChange))
	sb.WriteString(msg.Sprintf("BTC Dominance: *%.1f%% (%+.2f%%)*\n", data.BTCDom, data.BTCDom24HChange))

	wrUpDownEmj(&sb, float64(data.ETHDom24HChange))
	sb.WriteString(msg.Sprintf("ETH Dominance: *%.1f%% (%+.2f%%)*\n\n", data.ETHDom, data.ETHDom24HChange))

	// Time
	sb.WriteString(msg.Sprintf("Last update: _%s_\n", time.Unix(data.Upd, 0).UTC().Format(TIME_FORMAT)))

	return sb.String()
}

// Output must be send with tg.ModeMarkdown
func fmtFuturesStats(data models.Data) string {
	var sb strings.Builder
	sb.Grow(256)
	msg := message.NewPrinter(language.Ukrainian)

	// Header
	sb.WriteString("*Futures Stats:*\n\n")

	// Data
	wrUpDownEmj(&sb, float64(data.OpenInterest24HChange))
	sb.WriteString(msg.Sprintf("Open Interest: *$%d (%+.2f%%)*\n", data.OpenInterest, data.OpenInterest24HChange))

	wrUpDownEmj(&sb, float64(data.Liquidations24HUsdChange))
	sb.WriteString(msg.Sprintf("Liquidations 24H: *$%.f (%+.2f%%)*\n", data.Liquidation24HUsd, data.Liquidations24HUsdChange))

	sb.WriteString(msg.Sprintf("Number of liquidations: *%d*\n\n", data.Liquidation24HNum))

	// Time
	sb.WriteString(msg.Sprintf("Last update: _%s_\n", time.Unix(data.Upd, 0).UTC().Format(TIME_FORMAT)))

	return sb.String()
}

// Writes EMJ_GRAPHUP if val > 0, EMJ_GRAPHDOWN if val < 0, and empty string otherwise
func wrUpDownEmj(sb *strings.Builder, val float64) {
	if val > 0 {
		sb.WriteString(EMJ_GRAPHUP)
	} else if val < 0 {
		sb.WriteString(EMJ_GRAPHDOWN)
	}
}
