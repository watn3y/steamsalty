package config

type config struct {
	LogLevel int `env:"STEAMSALTY_LOGLEVEL, default=1"`

	ChatID        int64    `env:"STEAMSALTY_CHATID, required`
	Watchers      []uint64 `env:"STEAMSALTY_WATCHERS, required`
	SleepInterval int      `env:"STEAMSALTY_SLEEPINTERVAL, default=60"`

	Translate struct {
		Enabled  bool   `env:"STEAMSALTY_TRANSLATE_ENABLED, default=False"`
		Language string `env:"STEAMSALTY_TRANSLATE_LANGUAGE, default=EN-US"`
	}

	Telegram struct {
		APIToken string `env:"STEAMSALTY_TELEGRAM_APITOKEN, required"`
	}

	Steam struct {
		APIKey string `env:"STEAMSALTY_STEAMAPIKEY, required`
	}

	DeepL struct {
		APIKey   string `env:"STEAMSALTY_DEEPL_APIKEY"`
		FreeTier bool   `env:"STEAMSALTY_DEEPL_FREETIER, default=True"`
	}
}
