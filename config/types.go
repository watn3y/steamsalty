package config

type config struct {
	TelegramAPIToken string   `env:"TELEGRAMAPITOKEN, required"`
	SteamAPIKey      string   `env:"STEAMAPIKEY, required"`
	DebugMode        bool     `env:"DEBUGMODE, default=false"`
	ChatID           int64    `env:"CHATID"`
	Watchers         []uint64 `env:"WATCHERS"`
	SleepInterval    int      `env:"SLEEPINTERVAL"`
}
