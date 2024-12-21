package steam

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

func getComments(steamID uint64, start int, count int) (page CommentsPage) {

	baseURL := "https://steamcommunity.com/comment/Profile/render/"

	url, err := url.Parse(baseURL + strconv.FormatUint(steamID, 10))
	if err != nil {
		log.Error().Err(err).Msg("Unable to Parse SteamID into URL")
		return
	}

	query := url.Query()
	query.Set("start", strconv.Itoa(start))
	query.Set("count", strconv.Itoa(count))
	url.RawQuery = query.Encode()

	resp, err := http.Get(url.String())
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error().Err(err).Int("Response Code", resp.StatusCode).Msg("Failed to get Comments")
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse Comments")
		log.Trace().Interface("Body", resp.Body)
	}

	err = json.Unmarshal(body, &page)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse Comments as JSON")
	}

	log.Debug().Interface("CommentPage", page).Uint64("ProfileID", steamID).Msg("Successfully got Comment Page")

	return page
}

func parseComments(rawComments CommentsPage) (comments []Comment) {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawComments.CommentsHTML))
	if err != nil {
		log.Error().Err(err).Msg("Error while parsing CommentsHTML")
		return
	}
	doc.Find(".commentthread_comment.responsive_body_text").Each(func(i int, s *goquery.Selection) {
		var c Comment

		parsedID, err := strconv.ParseUint(strings.TrimPrefix(s.AttrOr("id", ""), "comment_"), 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("Error while parsing Comment ID")
			return
		}
		c.ID = parsedID

		c.Timestamp, err = strconv.ParseInt(s.Find(".commentthread_comment_timestamp").AttrOr("data-timestamp", ""), 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("Error while parsing Comment Timestamp")
			return
		}

		c.Author = s.Find(".commentthread_comment_author .hoverunderline bdi").Text()

		c.AuthorProfileURL, _ = s.Find(".commentthread_comment_author .hoverunderline").Attr("href")

		c.Text = strings.TrimSpace(s.Find(".commentthread_comment_text").Text())

		comments = append(comments, c)
	})

	slices.Reverse(comments)
	log.Debug().Interface("Comments", comments).Msg("Successfully parsed Comment Page")
	return comments
}
