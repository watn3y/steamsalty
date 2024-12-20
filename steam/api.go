package steam

import (
	"github.com/Philipp15b/go-steamapi"
	"github.com/rs/zerolog/log"
	"watn3y/steamsalty/config"
)

func getPlayerDetails(steamID uint64) (summary steamapi.PlayerSummary) {

	response, err := steamapi.GetPlayerSummaries([]uint64{steamID}, config.BotConfig.SteamAPIKey)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get Player Summary")
	}

	return response[0]
}