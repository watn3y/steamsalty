package main

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	
	"watn3y/steamsalty/config"
)

func main() {
	println("Starting SteamSalty...")

	configureLogger()

	config.LoadConfig()

	bot()

}

func configureLogger() {

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		const prefix = "steamsalty/"

		index := strings.Index(file, prefix)
		if index != -1 {
			return file[index+len(prefix):] + ":" + strconv.Itoa(line)
		}
		return file + ":" + strconv.Itoa(line)
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}

	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
	log.Info().Msg("Logger started successfully")

}
