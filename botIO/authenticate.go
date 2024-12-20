package botIO

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"watn3y/steamsalty/config"
)

func Authenticate() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(config.BotConfig.TelegramAPIToken)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to authenticate")
	}

	bot.Debug = config.BotConfig.DebugMode

	log.Info().Int64("ID", bot.Self.ID).Str("username", bot.Self.UserName).Msg("Successfully authenticated to Telegram API")

	return bot
}
