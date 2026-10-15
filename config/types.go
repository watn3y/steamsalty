package config

// Every env var is named after its field's path: STEAMSALTY_<SECTION>_<FIELD>.
// config_test.go enforces this.
type config struct {
	LogLevel int `env:"STEAMSALTY_LOGLEVEL, default=1"`

	Telegram struct {
		APIToken string `env:"STEAMSALTY_TELEGRAM_APITOKEN, required"`
		ChatID   int64  `env:"STEAMSALTY_TELEGRAM_CHATID, required"`
	}

	Steam struct {
		APIKey          string   `env:"STEAMSALTY_STEAM_APIKEY, required"`
		WatchedProfiles []uint64 `env:"STEAMSALTY_STEAM_WATCHEDPROFILES, required"`
		SleepInterval   int      `env:"STEAMSALTY_STEAM_SLEEPINTERVAL, default=60"`
	}

	DeepL struct {
		Enabled        bool   `env:"STEAMSALTY_DEEPL_ENABLED, default=False"`
		TargetLanguage string `env:"STEAMSALTY_DEEPL_TARGETLANGUAGE, default=EN-US"`
		APIKey         string `env:"STEAMSALTY_DEEPL_APIKEY"`
		FreeTier       bool   `env:"STEAMSALTY_DEEPL_FREETIER, default=True"`
	}
}
