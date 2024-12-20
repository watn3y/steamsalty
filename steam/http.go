package steam

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func getComments(steamID uint64, start int, count int) (comments CommentResponse) {

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

	err = json.Unmarshal(body, &comments)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse Comments as JSON")
	}

	log.Debug().Interface("CommentPage", comments).Msg("Successfully got Comment Page")
	return comments
}
