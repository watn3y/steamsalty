package steam

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

// httpClient is used for steamcommunity.com requests. The timeout keeps a hung
// connection from blocking a watcher forever.
var httpClient = &http.Client{Timeout: 30 * time.Second}

// GetComments fetches a page of profile comments. Every failure (transport error,
// non-200 status, unparsable body, success=false) is returned as an error, so a
// failed request can never be mistaken for a comment page that changed.
func GetComments(steamID uint64, start int, count int) (page CommentsPage, err error) {

	baseURL := "https://steamcommunity.com/comment/Profile/render/"

	url, err := url.Parse(baseURL + strconv.FormatUint(steamID, 10))
	if err != nil {
		return CommentsPage{}, fmt.Errorf("build comments URL: %w", err)
	}

	query := url.Query()
	query.Set("start", strconv.Itoa(start))
	query.Set("count", strconv.Itoa(count))
	url.RawQuery = query.Encode()

	resp, err := httpClient.Get(url.String())
	if err != nil {
		return CommentsPage{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CommentsPage{}, fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CommentsPage{}, fmt.Errorf("read comments response: %w", err)
	}

	err = json.Unmarshal(body, &page)
	if err != nil {
		return CommentsPage{}, fmt.Errorf("parse comments as JSON: %w", err)
	}

	if !page.Success {
		return CommentsPage{}, errors.New("steam returned success=false")
	}

	log.Trace().Interface("CommentPage", page).Uint64("ProfileID", steamID).Msg("Retrieved Comment Page successfully")

	return page, nil
}

func parseComments(rawComments CommentsPage) (comments []Comment) {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawComments.CommentsHTML))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse CommentsHTML")
		return
	}
	doc.Find(".commentthread_comment.responsive_body_text").Each(func(i int, s *goquery.Selection) {
		var c Comment

		parsedID, err := strconv.ParseUint(strings.TrimPrefix(s.AttrOr("id", ""), "comment_"), 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse Comment ID")
			return
		}
		c.ID = parsedID

		// Steam renders an empty .commentthread_comment_timestamp element before the
		// one carrying data-timestamp, so select by the attribute.
		c.Timestamp, err = strconv.ParseInt(s.Find(".commentthread_comment_timestamp[data-timestamp]").AttrOr("data-timestamp", ""), 10, 64)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse Comment Timestamp")
			return
		}

		c.Author = s.Find(".commentthread_comment_author .hoverunderline bdi").Text()

		c.AuthorProfileURL, _ = s.Find(".commentthread_comment_author .hoverunderline").Attr("href")

		c.Text = strings.TrimSpace(s.Find(".commentthread_comment_text").Text())

		comments = append(comments, c)
	})

	slices.Reverse(comments)
	log.Trace().Interface("Comments", comments).Msg("Parsed comment page successfully")
	return comments
}
