package steam

import (
	"errors"
	"fmt"
	"net/url"
	"watn3y/steamsalty/config"

	"github.com/Philipp15b/go-steamapi"
	"github.com/rs/zerolog/log"
)

func GetPlayerDetails(steamID uint64) (summary steamapi.PlayerSummary, err error) {

	response, err := steamapi.GetPlayerSummaries([]uint64{steamID}, config.BotConfig.Steam.APIKey)
	if err != nil {
		// Transport errors are *url.Error values whose message contains the full
		// request URL, including the key= query parameter. Drop the URL so the API
		// key never ends up in the logs.
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			err = urlErr.Err
		}
		return summary, fmt.Errorf("request player summary: %w", err)
	}
	if len(response) == 0 {
		return summary, fmt.Errorf("steam API returned no player summary for SteamID %d", steamID)
	}

	log.Debug().Interface("Player", response[0]).Msg("Retrived player summary from SteamAPI successfully")
	return response[0], nil
}
