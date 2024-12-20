package config

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	envconfig "github.com/sethvargo/go-envconfig"
)

var BotConfig config

func LoadConfig() {
	if err := envconfig.Process(context.Background(), &BotConfig); err != nil {
		log.Panic().Err(err).Msg("error parsing config from env variables")
	}

	if !BotConfig.DebugMode {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Info().Msg("Loaded config")
	log.Debug().Interface("config", BotConfig).Msg("")

}
