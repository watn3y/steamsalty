package config

type config struct {
	LogLevel int `env:"STEAMSALTY_LOGLEVEL" default:"1"`

	SleepInterval     int      `env:"STEAMSALTY_SLEEPINTERVAL" default:"60"`
	ChatID            int64    `env:"STEAMSALTY_CHATID" required:"true"`
	Watchers          []uint64 `env:"STEAMSALTY_WATCHERS" required:"true"`
	TranslateEnabled  bool     `env:"STEAMSALTY_TRANSLATE_ENABLED" default:"false"`
	TranslateLanguage string   `env:"STEAMSALTY_TRANSLATE_LANGUAGE" default:"EN-US"`
	TelegramAPIToken  string   `env:"STEAMSALTY_TELEGRAMAPITOKEN" required:"true"`
	SteamAPIKey       string   `env:"STEAMSALTY_STEAMAPIKEY" required:"true"`
	
	DeepL struct {
		APIKey    string `env:"APIKEY"`
		FreeTier  bool   `env:"FREETIER" default:"true"`
	} `env:"STEAMSALTY_DEEPL_"`
}

