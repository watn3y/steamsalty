package steam

import (
	"watn3y/steamsalty/config"

	"github.com/Philipp15b/go-steamapi"
	"github.com/rs/zerolog/log"
)

func GetPlayerDetails(steamID uint64) (summary steamapi.PlayerSummary) {

	response, err := steamapi.GetPlayerSummaries([]uint64{steamID}, config.BotConfig.SteamAPIKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get Player Summary")
	}
	log.Debug().Interface("Player", response[0]).Msg("Got PlayerSummary from Steam API")
	return response[0]
}
