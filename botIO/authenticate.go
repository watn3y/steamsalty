package botIO

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"watn3y/steamsalty/config"
)

func Authenticate() (tgbotapi.UpdatesChannel, *tgbotapi.BotAPI) {
	bot, err := tgbotapi.NewBotAPI(config.BotConfig.TelegramAPIToken)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to authenticate")
	}

	bot.Debug = config.BotConfig.DebugMode

	updates := tgbotapi.NewUpdate(0)
	updates.Timeout = 60

	log.Info().Int64("ID", bot.Self.ID).Str("username", bot.Self.UserName).Msg("Successfully authenticated to Telegram API")

	return bot.GetUpdatesChan(updates), bot

}
