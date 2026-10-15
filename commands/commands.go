package commands

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"watn3y/steamsalty/botIO"
	"watn3y/steamsalty/config"
	"watn3y/steamsalty/steam"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func SetBotCommands(bot *tgbotapi.BotAPI) {
	github := tgbotapi.BotCommand{Command: "github", Description: "Source GitHub repo"}
	info := tgbotapi.BotCommand{Command: "info", Description: "Summary of watched profiles"}

	commands := tgbotapi.NewSetMyCommands(github, info)

	result, err := bot.Request(commands)
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish commands to Telegram")
		return
	}

	log.Debug().Interface("commands", result).Msg("Published commands to Telegram successfully")
}

func Commands(update tgbotapi.Update, bot *tgbotapi.BotAPI) {

	cmd := strings.ToLower(update.Message.Command())

	switch cmd {
	case "start":
		startGithub(update, bot)
	case "github":
		startGithub(update, bot)
	case "info":
		info(update, bot)
	}
}

func startGithub(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	message := tgbotapi.MessageConfig{
		BaseChat:              tgbotapi.BaseChat{ChatID: update.Message.Chat.ID, ReplyToMessageID: update.Message.MessageID},
		ParseMode:             "html",
		DisableWebPagePreview: false,
		Text:                  "Check out: https://github.com/watn3y/steamsalty",
	}
	botIO.SendMessage(message, bot)
}

func info(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	if update.Message.Chat.ID != config.BotConfig.ChatID {
		return
	}

	textInfo := "<b>Watching profiles:</b> " + "\n"

	for _, steamID := range config.BotConfig.Watchers {
		profile, err := steam.GetPlayerDetails(steamID)
		if err != nil {
			log.Error().Err(err).Uint64("ProfileID", steamID).Msg("Failed to retrieve player summary, using SteamID")
			profile.PersonaName = strconv.FormatUint(steamID, 10)
			profile.ProfileURL = "https://steamcommunity.com/profiles/" + strconv.FormatUint(steamID, 10) + "/"
		}

		lastComment := "never :("
		commentCount := "unknown"
		comments, err := steam.GetComments(steamID, 0, 0)
		if err != nil {
			log.Error().Err(err).Uint64("ProfileID", steamID).Msg("Failed to retrieve comments")
			lastComment = "unknown (failed to reach Steam)"
		} else {
			commentCount = strconv.Itoa(comments.TotalCount)
			if comments.TimeLastPost > 0 {
				lastComment = time.Unix(comments.TimeLastPost, 0).Format(time.RFC1123)
			}
		}

		textInfo += fmt.Sprintf(`<b><a href="%s">%s</a></b>:`, html.EscapeString(profile.ProfileURL), html.EscapeString(profile.PersonaName)) + "\n" +
			fmt.Sprintf(`Last Comment: %s`, lastComment) + "\n" +
			fmt.Sprintf(`Number of Comments: %s`, commentCount) + "\n\n"
	}

	message := tgbotapi.MessageConfig{
		BaseChat:              tgbotapi.BaseChat{ChatID: update.Message.Chat.ID, ReplyToMessageID: update.Message.MessageID},
		ParseMode:             "html",
		DisableWebPagePreview: true,
		Text:                  textInfo,
	}
	botIO.SendMessage(message, bot)

}
