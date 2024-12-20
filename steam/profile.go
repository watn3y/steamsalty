package steam

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"math"
	"sync"
	"time"
	"watn3y/steamsalty/botIO"
	"watn3y/steamsalty/config"
)

var sleeptime time.Duration = 10 * time.Second

func StartWatchers(bot *tgbotapi.BotAPI) {
	var wg sync.WaitGroup

	for _, steamID := range config.BotConfig.Watchers {
		wg.Add(1)
		go func(steamID uint64) {
			defer wg.Done()
			watcher(bot, steamID)
		}(steamID)
	}

	wg.Wait()
}

func watcher(bot *tgbotapi.BotAPI, steamID uint64) {
	log.Info().Uint64("SteamID", steamID).Msg("Started Watcher")
	var previousCount int
	for {
		currentCount := getComments(steamID, math.MaxInt32, 0).TotalCount
		if previousCount == 0 || currentCount <= previousCount {
			previousCount = currentCount
			time.Sleep(sleeptime)
			continue
		}

		log.Info().Int("NumComment", currentCount).Uint64("SteamID", steamID).Msg("Found new comment")

		player := getPlayerDetails(steamID)

		msg := tgbotapi.MessageConfig{
			BaseChat:              tgbotapi.BaseChat{ChatID: config.BotConfig.ChatID},
			ParseMode:             "html",
			DisableWebPagePreview: false,
			Text:                  fmt.Sprintf(`New comment on <a href="%s">%s's</a> profile`, player.ProfileURL, player.PersonaName),
		}

		botIO.SendMessage(msg, bot)

		previousCount = currentCount
		time.Sleep(sleeptime)
	}
}
