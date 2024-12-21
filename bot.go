package main

import (
	"time"

	"watn3y/steamsalty/botIO"
	"watn3y/steamsalty/commands"
	"watn3y/steamsalty/steam"

	"github.com/rs/zerolog/log"
)

func bot() {
	updates, bot := botIO.Authenticate()

	go steam.StartWatchers(bot)

	for update := range updates {
		log.Debug().Interface("update", update).Msg("Received update")

		if update.Message == nil || update.Message.Text == "" {
			log.Debug().Int("UpdateID", update.UpdateID).Msg("Unable to parse update")
			continue
		}
		if update.Message.Time().UTC().Unix() < time.Now().UTC().Unix() {
			log.Debug().Int("UpdateID", update.UpdateID).Msg("Skipping old update")
			continue
		}

		log.Info().Int64("ChatID", update.Message.Chat.ID).Int64("UserID", update.Message.From.ID).Str("Text", update.Message.Text).Msg("Recieved Message")

		if update.Message.IsCommand() {
			commands.Commands(update, bot)
		}
	}
}
