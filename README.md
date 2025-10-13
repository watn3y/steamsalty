# SteamSalty

SteamSalty notifies you on telegram about new comments on any steam profile **with built in auto translation**.

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
      - STEAMSALTY_CHATID=
      - STEAMSALTY_WATCHERS=
      - STEAMSALTY_SLEEPINTERVAL=
      #- STEAMSALTY_TRANSLATE_ENABLED=
      #- STEAMSALTY_TRANSLATE_LANGUAGE=
      - STEAMSALTY_TELEGRAM_APITOKEN=
      - STEAMSALTY_STEAMAPIKEY=
      #- STEAMSALTY_DEEPL_APIKEY=
      #- STEAMSALTY_DEEPL_FREETIER=
```

## Running on Linux

Grab a release from the [releases page](https://github.com/watn3y/steamsalty/releases). Make sure to set your  [environment variables](#environment-variables) accordingly.

## Environment Variables

> [!NOTE]  
> For development purposes, SteamSalty supports loading environment variables from a .env file placed in the project root directory.

| Variable                        | Description                                                                                                                                                             | Default  | Required | Example                                       |
|---------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|----------|-----------------------------------------------|
| `STEAMSALTY_LOGLEVEL`           | LogLevel as described [in the zerolog documentation](https://pkg.go.dev/github.com/rs/zerolog@v1.34.0#readme-simple-leveled-logging-example)                            | 1 (Info) | ❌        | 1                                             |
| `STEAMSALTY_SLEEPINTERVAL`      | Amount of time to wait between requests to Steam in seconds                                                                                                             | 60       | ❌        | 60                                            |
| `STEAMSALTY_TRANSLATE_ENABLED`  | Whether to enable translation of comments                                                                                                                               | False    | ❌        | True, False                                   |
| `STEAMSALTY_TRANSLATE_LANGUAGE` | Language to translate as described [in the DeepL API documentation](https://developers.deepl.com/docs/getting-started/supported-languages#translation-target-languages) | EN-US    | ❌        | EN-US,DE                                      |
| `STEAMSALTY_CHATID`             | Chat to notify about new comments                                                                                                                                       | None     | ✅        | -1001234567890                                |
| `STEAMSALTY_WATCHERS`           | SteamIDs (in SteamID64 format) to check for new profile comments                                                                                                        | None     | ✅        | 76561198012345678,76561198087654321           |
| `STEAMSALTY_TELEGRAMAPITOKEN`   | Telegram BotToken, get it from  [@BotFather on Telegram](https://t.me/BotFather)                                                                                        | None     | ✅        | 1234567890:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw |
| `STEAMSALTY_STEAMAPIKEY`        | Steam API Key, get it from [steamcommunity.com/dev/apikey](https://steamcommunity.com/dev/apikey)                                                                       | None     | ✅        | A7B3C9D2E5F1A4B8C6D9E2F5A8B1C4D7E0F3A6B9      |
| `STEAMSALTY_DEEPL_APIKEY`       | DeepL API Key, get it from [deepl.com/en/your-account/keys](https://www.deepl.com/en/your-account/keys)                                                                 | None     | ❌        | a1b2c3d4-56e7-89f0-a1b2-c3d4e5f6a7b8:fx       |
| `STEAMSALTY_DEEPL_FREETIER`     | Whether you are using the DeepL Free Tier                                                                                                                               | True     | ❌        | True, False                                   |

## Nice to know

### Semantic Versioning

This project does it's best to follow [Semantic Versioning](https://semver.org/#semantic-versioning-200), however I can't guarantee anything.
