package config

type config struct {
	LogLevel         int      `env:"STEAMSALTY_LOGLEVEL, default=1"`
	TelegramAPIToken string   `env:"STEAMSALTY_TELEGRAMAPITOKEN, required"`
	SteamAPIKey      string   `env:"STEAMSALTY_STEAMAPIKEY, required"`
	ChatID           int64    `env:"STEAMSALTY_CHATID, required"`
	Watchers         []uint64 `env:"STEAMSALTY_WATCHERS, required"`
	SleepInterval    int      `env:"STEAMSALTY_SLEEPINTERVAL, default=60"`
}
