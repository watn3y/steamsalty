# SteamSalty

SteamSalty notifies you on telegram about new comments on any steam profile **with built in auto translation**.

> [!IMPORTANT]  
> The version of this repository you are seeing is likely a mirror.
>
> If you want to contribute any changes please head over to [git.watn3y.de](https://git.watn3y.de/watn3y/steamsalty) and register an account.

## Running with Docker Compose

The Docker image is available on these registries.

- [git.watn3y.de](https://git.watn3y.de/watn3y/-/packages/container/steamsalty/versions) - includes all tags
- [Docker Hub](https://hub.docker.com/r/watn3y/steamsalty) - includes only releases

Example compose file:

```yaml
services:
  steamsalty:
    image: watn3y/steamsalty:v1.0.0
    # image: git.watn3y.de/watn3y/steamsalty:latest # Use :latest to be up to date with the master branch on git
    container_name: steamsalty
    restart: unless-stopped
    volumes:
      - /etc/localtime:/etc/localtime:ro
    environment:
      #STEAMSALTY_LOGLEVEL: 
      STEAMSALTY_TELEGRAM_APITOKEN: 
      STEAMSALTY_TELEGRAM_CHATID: 
      STEAMSALTY_STEAM_APIKEY:
      STEAMSALTY_STEAM_WATCHEDPROFILES:
      #STEAMSALTY_STEAM_SLEEPINTERVAL:
      #STEAMSALTY_DEEPL_ENABLED:
      #STEAMSALTY_DEEPL_TARGETLANGUAGE:
      #STEAMSALTY_DEEPL_APIKEY:
      #STEAMSALTY_DEEPL_FREETIER:
```

## Running on Linux

Grab a release from the [releases page](https://git.watn3y.de/watn3y/steamsalty/releases). Make sure to set your  [environment variables](#environment-variables) accordingly.

## Environment Variables

> [!NOTE]  
> For development purposes, SteamSalty supports loading environment variables from a .env file placed in the project root directory.
>
> Optional variables that are unset or empty use their default.

| Variable                           | Description                                                                                                                                                                 | Default  | Required | Example                                       |
|------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|----------|-----------------------------------------------|
| `STEAMSALTY_LOGLEVEL`              | LogLevel as described [in the zerolog documentation](https://pkg.go.dev/github.com/rs/zerolog@v1.34.0#readme-simple-leveled-logging-example)                                | 1 (Info) | ❌        | 1                                             |
| `STEAMSALTY_TELEGRAM_APITOKEN`     | Telegram BotToken, get it from [@BotFather on Telegram](https://t.me/BotFather)                                                                                             | None     | ✅        | 1234567890:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw |
| `STEAMSALTY_TELEGRAM_CHATID`       | Chat to notify about new comments                                                                                                                                           | None     | ✅        | -1001234567890                                |
| `STEAMSALTY_STEAM_APIKEY`          | Steam API Key, get it from [steamcommunity.com/dev/apikey](https://steamcommunity.com/dev/apikey)                                                                           | None     | ✅        | A7B3C9D2E5F1A4B8C6D9E2F5A8B1C4D7E0F3A6B9      |
| `STEAMSALTY_STEAM_WATCHEDPROFILES` | SteamIDs (in SteamID64 format) to check for new profile comments                                                                                                            | None     | ✅        | 76561198012345678,76561198087654321           |
| `STEAMSALTY_STEAM_SLEEPINTERVAL`   | Amount of time to wait between requests to Steam in seconds, at least 1                                                                                                     | 60       | ❌        | 60                                            |
| `STEAMSALTY_DEEPL_ENABLED`         | Whether to enable translation of comments                                                                                                                                   | False    | ❌        | True, False                                   |
| `STEAMSALTY_DEEPL_TARGETLANGUAGE`  | Language to translate to, as described [in the DeepL API documentation](https://developers.deepl.com/docs/getting-started/supported-languages#translation-target-languages) | EN-US    | ❌        | EN-US, DE                                     |
| `STEAMSALTY_DEEPL_APIKEY`          | DeepL API Key, get it from [deepl.com/en/your-account/keys](https://www.deepl.com/en/your-account/keys). Required if `STEAMSALTY_DEEPL_ENABLED` is True                     | None     | ❌        | a1b2c3d4-56e7-89f0-a1b2-c3d4e5f6a7b8:fx       |
| `STEAMSALTY_DEEPL_FREETIER`        | Whether you are using the DeepL Free Tier                                                                                                                                   | True     | ❌        | True, False                                   |

## Nice to know

### Semantic Versioning

This project does it's best to follow [Semantic Versioning](https://semver.org/#semantic-versioning-200), however I can't guarantee anything.

### AI Disclosure

AI was used in the process of creating this project. Mainly for the Go unit tests.
