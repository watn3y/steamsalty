package config

import (
	"context"
	"fmt"
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

	cfg, err := load(envconfig.OsLookuper())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config from env variables")
	}
	BotConfig = cfg
	zerolog.SetGlobalLevel(zerolog.Level(BotConfig.LogLevel))

	log.Info().Msg("Config loaded successfully")
	log.Debug().Interface("config", BotConfig).Msg("")

}


func load(env envconfig.Lookuper) (config, error) {
	var c config
	err := envconfig.ProcessWith(context.Background(), &envconfig.Config{
		Target:   &c,
		Lookuper: nonEmptyLookuper{env},
	})
	if err != nil {
		return config{}, err
	}

	c.DeepL.TargetLanguage = strings.ToUpper(c.DeepL.TargetLanguage)


	if c.Steam.SleepInterval < 1 {
		return config{}, fmt.Errorf("STEAMSALTY_STEAM_SLEEPINTERVAL must be at least 1 second, got %d", c.Steam.SleepInterval)
	}

	return c, nil
}


type nonEmptyLookuper struct{ envconfig.Lookuper }

func (l nonEmptyLookuper) Lookup(key string) (string, bool) {
	value, found := l.Lookuper.Lookup(key)
	return value, found && value != ""
}
