# SteamSalty

SteamSalty notifies you on telegram about new comments on any steam profile.

## Running with Docker Compose

Docker image: <https://hub.docker.com/r/watn3y/steamsalty>

Example compose file:

```yaml
services:
  steamsalty:
    image: watn3y/steamsalty:latest # use :<branchname> to be up-to-date with any branch 
    container_name: steamsalty
    restart: unless-stopped
    volumes:
      - /etc/localtime:/etc/localtime:ro
    environment:
      #- STEAMSALTY_LOGLEVEL=
      - STEAMSALTY_TELEGRAMAPITOKEN=
      - STEAMSALTY_STEAMAPIKEY=
      - STEAMSALTY_CHATID=
      - STEAMSALTY_WATCHERS=
      #- STEAMSALTY_SLEEPINTERVAL=
```

## Running on Linux

Grab a release from the [releases page](https://github.com/watn3y/steamsalty/releases). Make sure to set your **environment variables** accordingly.

## Environment Variables

| Variable                      | Description                                                                                                                                   | Default            |
| ----------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------| ------------------ |
| `STEAMSALTY_LOGLEVEL`         | LogLevel as described [in the zerolog documentation](https://pkg.go.dev/github.com/rs/zerolog@v1.34.0#readme-simple-leveled-logging-example)  | 1 (Info)           |
| `STEAMSALTY_TELEGRAMAPITOKEN` | Telegram Bot Token, get from @BotFather                                                                                                       | None, **required** |
| `STEAMSALTY_STEAMAPIKEY`      | Steam API Key, get from [steamcommunity.com/dev/apikey](https://steamcommunity.com/dev/apikey)                                                | None, **required** |
| `STEAMSALTY_CHATID`           | Chat to notify about new Comments                                                                                                             | None, **required** |
| `STEAMSALTY_WATCHERS`         | SteamIDs (in SteamID64 format) to check for new Profile Comments                                                                              | None, **required** |
| `STEAMSALTY_SLEEPINTERVAL`    | Amount of time to wait between requests to Steam in seconds                                                                                   | 60                 |
