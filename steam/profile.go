package steam

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
	"watn3y/steamsalty/botIO"
	"watn3y/steamsalty/config"
	deepl "watn3y/steamsalty/deepL"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const steamContentCheckText string = "This comment is awaiting analysis by our automated content check system. It will be temporarily hidden until we verify that it does not contain harmful content (e.g. links to websites that attempt to steal information)."

func StartWatchers(bot *tgbotapi.BotAPI) {

	var wg sync.WaitGroup

	for _, steamID := range config.BotConfig.Watchers {
		wg.Add(1)
		go func(steamID uint64) {
			defer wg.Done()
			watcher(bot, steamID, time.Duration(config.BotConfig.SleepInterval)*time.Second)
		}(steamID)
	}

	wg.Wait()
}

func watcher(bot *tgbotapi.BotAPI, steamID uint64, sleeptime time.Duration) {
	log.Info().Uint64("SteamID", steamID).Msg("Starting Watcher")

	var newestProcessedComment int64 = 0

	for {
		currentCommentsPage := GetComments(steamID, 0, math.MaxInt32)
		if newestProcessedComment == 0 || newestProcessedComment == currentCommentsPage.TimeLastPost {
			newestProcessedComment = currentCommentsPage.TimeLastPost
			time.Sleep(sleeptime)
			continue
		}

		if strings.Contains(currentCommentsPage.CommentsHTML, steamContentCheckText) {
			log.Info().Uint64("ProfileID", steamID).Msg("Found new comment(s) still being checked by Steam, retrying in " + sleeptime.String())
			time.Sleep(sleeptime)
			continue
		}

		log.Info().Uint64("ProfileID", steamID).Msg("Found new comment(s)")

		profileOwner := GetPlayerDetails(steamID)

		for _, comment := range parseComments(currentCommentsPage) {
			log.Debug().Interface("Comment", comment).Msg("Processing comment")
			if comment.Timestamp <= newestProcessedComment {
				log.Debug().Uint64("CommentID", comment.ID).Msg("Skipping comment")
				continue
			}

			msg := tgbotapi.MessageConfig{
				BaseChat:              tgbotapi.BaseChat{ChatID: config.BotConfig.ChatID},
				ParseMode:             "HTML",
				DisableWebPagePreview: true,
				Text: fmt.Sprintf(`<b><a href="%s">%s</a> just commented on <a href="%s">%s</a>'s profile:</b>`, comment.AuthorProfileURL, comment.Author, profileOwner.ProfileURL, profileOwner.PersonaName) + "\n" +
					"<blockquote>" + comment.Text + "</blockquote>",
			}

			if config.BotConfig.Translate.Enabled {
				translatedText, translatedTextLanguage, err := deepl.Translate(comment.Text)
				if err != nil {
					log.Error().Err(err).Msg("Failed to translate comment, continuing without translation")
					msg.Text += "\n" + "Translation failed"

				} else if !languageMatches(translatedTextLanguage, config.BotConfig.Translate.Language) {
					msg.Text += "\n" +
						fmt.Sprintf(`Translated from %s:`, deepl.SourceLanguages[translatedTextLanguage]) + "\n" +
						"<blockquote>" + translatedText + "</blockquote>"
				}
			}

			log.Info().Interface("Comment", comment).Msg("Notifying about new comment")
			botIO.SendMessage(msg, bot)
			time.Sleep(time.Second * 3) //I have no Idea why this is here
		}

		newestProcessedComment = currentCommentsPage.TimeLastPost
		time.Sleep(sleeptime)
	}
}

func languageMatches(sourceLanguage, targetLang string) (languageMatches bool) {
	apiUpper := strings.ToUpper(sourceLanguage)
	configUpper := strings.ToUpper(targetLang)

	if apiUpper == configUpper {
		return true
	}

	if strings.HasPrefix(configUpper, apiUpper+"-") {
		return true
	}

	return false
}
