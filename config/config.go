package config

import (
	"context"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	envconfig "github.com/sethvargo/go-envconfig"
)

var BotConfig config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Info().Err(err).Msg("Failed to load .env file, using the system environment")
	} else {
		log.Info().Err(err).Msg(".env file loaded successfully")
	}

	if err := envconfig.Process(context.Background(), &BotConfig); err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config from env variables")
	}
	zerolog.SetGlobalLevel(zerolog.Level(BotConfig.LogLevel))

	BotConfig.Translate.Language = strings.ToUpper(BotConfig.Translate.Language)
	

	log.Info().Msg("Config loaded successfully")
	log.Debug().Interface("config", BotConfig).Msg("")

}
