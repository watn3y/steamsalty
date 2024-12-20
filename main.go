package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"watn3y/steamsalty/config"
)

func main() {
	fmt.Println("Starting SteamSalty...")

	configureLogger()

	config.LoadConfig()

	bot()

}

func configureLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}

	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()

	//! note that we overwrite the loglevel after loading the config in config/config.go. This is just the default
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	log.Info().Msg("Started Logger")

}
