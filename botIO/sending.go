package botIO

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

func SendMessage(message tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) (result tgbotapi.Message) {
	result, err := bot.Send(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message")
		return
	}

	log.Info().Int64("chat", result.Chat.ID).Str("msg", result.Text).Msg("Sent message successfully")
	log.Debug().Interface("msg", result).Msg("")

	return result
}

func EditMessage(message tgbotapi.EditMessageTextConfig, bot *tgbotapi.BotAPI) (result tgbotapi.Message) {
	result, err := bot.Send(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to edit message")
		return
	}

	log.Info().Int64("chat", result.Chat.ID).Str("msg", result.Text).Msg("Edited message successfully")
	log.Debug().Interface("msg", result).Msg("")

	return result
}

func SendVideo(message tgbotapi.VideoConfig, bot *tgbotapi.BotAPI) (result tgbotapi.Message) {
	result, err := bot.Send(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send video")
		return
	}

	log.Info().Int64("chat", result.Chat.ID).Msg("Sent video successfully")
	log.Debug().Interface("video", result).Msg("")

	return result
}

func SendPhoto(message tgbotapi.PhotoConfig, bot *tgbotapi.BotAPI) (result tgbotapi.Message) {
	result, err := bot.Send(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send photo")
		return
	}

	log.Info().Int64("chat", result.Chat.ID).Msg("Sent photo successfully")
	log.Debug().Interface("photo", result).Msg("")

	return result
}

func SendSticker(message tgbotapi.StickerConfig, bot *tgbotapi.BotAPI) (result tgbotapi.Message) {
	result, err := bot.Send(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send sticker")
		return
	}

	log.Info().Int64("chat", result.Chat.ID).Msg("Sent sticker successfully")
	log.Debug().Interface("sticker", result).Msg("")

	return result
}
