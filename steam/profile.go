package steam

import (
	"fmt"
	"html"
	"math"
	"strconv"
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

	for _, steamID := range config.BotConfig.Steam.WatchedProfiles {
		wg.Add(1)
		go func(steamID uint64) {
			defer wg.Done()
			watcher(bot, steamID, time.Duration(config.BotConfig.Steam.SleepInterval)*time.Second)
		}(steamID)
	}

	wg.Wait()
}

func watcher(bot *tgbotapi.BotAPI, steamID uint64, sleeptime time.Duration) {
	log.Info().Uint64("SteamID", steamID).Msg("Starting Watcher")

	var newestProcessedComment int64 = 0

	for {
		currentCommentsPage, err := GetComments(steamID, 0, math.MaxInt32)
		if err != nil {
			// A failed request says nothing about new comments. Leave newestProcessedComment
			// untouched so anything posted during the outage is picked up by the next successful poll.
			log.Warn().Err(err).Uint64("ProfileID", steamID).Msg("Failed to retrieve comments, retrying in " + sleeptime.String())
			time.Sleep(sleeptime)
			continue
		}

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

		profileOwner, err := GetPlayerDetails(steamID)
		if err != nil {
			// Still notify; name the profile by its SteamID instead of the persona name.
			log.Error().Err(err).Uint64("ProfileID", steamID).Msg("Failed to retrieve profile owner, using SteamID in notification")
			profileOwner.PersonaName = strconv.FormatUint(steamID, 10)
			profileOwner.ProfileURL = "https://steamcommunity.com/profiles/" + strconv.FormatUint(steamID, 10) + "/"
		}

		for _, comment := range parseComments(currentCommentsPage) {
			log.Debug().Interface("Comment", comment).Msg("Processing comment")
			if comment.Timestamp <= newestProcessedComment {
				log.Debug().Uint64("CommentID", comment.ID).Msg("Skipping comment")
				continue
			}

			msg := tgbotapi.MessageConfig{
				BaseChat:              tgbotapi.BaseChat{ChatID: config.BotConfig.Telegram.ChatID},
				ParseMode:             "HTML",
				DisableWebPagePreview: true,
				Text: fmt.Sprintf(`<b><a href="%s">%s</a> just commented on <a href="%s">%s</a>'s profile:</b>`,
					html.EscapeString(comment.AuthorProfileURL), html.EscapeString(comment.Author),
					html.EscapeString(profileOwner.ProfileURL), html.EscapeString(profileOwner.PersonaName)) + "\n" +
					"<blockquote>" + html.EscapeString(comment.Text) + "</blockquote>",
			}

			if config.BotConfig.DeepL.Enabled {
				translatedText, translatedTextLanguage, err := deepl.Translate(comment.Text)
				if err != nil {
					log.Error().Err(err).Msg("Failed to translate comment, continuing without translation")
					msg.Text += "\n" + "Translation failed"

				} else if !languageMatches(translatedTextLanguage, config.BotConfig.DeepL.TargetLanguage) && strings.TrimSpace(comment.Text) != strings.TrimSpace(translatedText) {
					msg.Text += "\n" +
						fmt.Sprintf(`Translated from %s:`, html.EscapeString(deepl.SourceLanguages[translatedTextLanguage])) + "\n" +
						"<blockquote>" + html.EscapeString(translatedText) + "</blockquote>"
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
