package commands

import (
	"fmt"
	"strings"
	"watn3y/steamsalty/botIO"
	"watn3y/steamsalty/config"
	"watn3y/steamsalty/steam"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func Commands(update tgbotapi.Update, bot *tgbotapi.BotAPI) {

	cmd := strings.ToLower(update.Message.Command())

	log.Debug().Str("cmd", cmd).Msg("Matching command")

	switch cmd {
	case "start":
		start(update, bot)
	case "info":
		info(update, bot)
	}
}

func start(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	message := tgbotapi.MessageConfig{
		BaseChat:              tgbotapi.BaseChat{ChatID: update.Message.Chat.ID, ReplyToMessageID: update.Message.MessageID},
		ParseMode:             "html",
		DisableWebPagePreview: false,
		Text:                  "https://github.com/watn3y/steamsalty",
	}
	botIO.SendMessage(message, bot)
}

func info(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message.Chat.ID != config.BotConfig.ChatID {
		return
	}

	textInfo := "<b>Watching profiles:</b> " + "\n"

	for _, steamID := range config.BotConfig.Watchers {
		profile := steam.GetPlayerDetails(steamID)

		textInfo += fmt.Sprintf(`- <a href="%s">%s</a>`, profile.ProfileURL, profile.PersonaName) + "\n"

	}

	message := tgbotapi.MessageConfig{
		BaseChat:              tgbotapi.BaseChat{ChatID: update.Message.Chat.ID, ReplyToMessageID: update.Message.MessageID},
		ParseMode:             "html",
		DisableWebPagePreview: true,
		Text:                  textInfo,
	}
	botIO.SendMessage(message, bot)

}
