package deepl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"watn3y/steamsalty/config"

	"github.com/rs/zerolog/log"
)

var client *apiClient

var SourceLanguages map[string]string
var TargetLanguages map[string]string

func Init() {
	log.Info().Msg("Translation is enabled, creating HTTP client for DeepL API")
	baseURL := baseURLPro
	if config.BotConfig.DeepLFreeTier {
		baseURL = baseURLFree
	}

	client = &apiClient{
		authKey:    config.BotConfig.DeepLAPIKey,
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	err := getAndValidateLanguages()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to set up languages")
	}
}

func Translate(text string) (translatedText string, sourceLanguage string, err error) {
	log.Debug().Str("text", text).Msg("Starting translation")
	req := translateRequest{
		Text:       []string{text},
		TargetLang: config.BotConfig.TranslateLanguage,
	}

	body, err := json.Marshal(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to ")
		return "", "", err
	}

	respBody, err := post("/translate", body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to make DeepL API request")
		return "", "", err
	}

	var result translateResponse
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse DeepL API response")
		return "", "", err
	}

	return result.Translations[0].Text, result.Translations[0].DetectedSourceLanguage, nil
}

func getAndValidateLanguages() (err error) {
	log.Info().Msg("Setting up supported languages")
	respBody, err := get("/languages?type=source")
	if err != nil {
		log.Error().Err(err).Msg("Failed to make DeepL API request")
		return err
	}

	var parsedResp languagesResponse
	err = json.Unmarshal(respBody, &parsedResp)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse DeepL API response")
		return err
	}

	SourceLanguages = make(map[string]string)
	for _, l := range parsedResp {
		SourceLanguages[l.Language] = l.Name
	}

	respBody, err = get("/languages?type=target")
	if err != nil {
		log.Error().Err(err).Msg("Failed to make DeepL API request")
		return err
	}

	err = json.Unmarshal(respBody, &parsedResp)
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse DeepL API response")
		return err
	}

	TargetLanguages = make(map[string]string)
	for _, l := range parsedResp {
		TargetLanguages[l.Language] = l.Name
	}

	if _, ok := TargetLanguages[config.BotConfig.TranslateLanguage]; !ok {
		return fmt.Errorf("Selected language not supported by DeepL")
	}

	return nil
}
