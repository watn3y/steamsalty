package deepl

import (
	"net/http"
)

type translateResponse struct {
	Translations []struct {
		DetectedSourceLanguage string `json:"detected_source_language"`
		Text                   string `json:"text"`
	} `json:"translations"`
}

type translateRequest struct {
	Text       []string `json:"text"`
	TargetLang string   `json:"target_lang"`
}

type languagesResponse []struct {
	Language string `json:"language"`
	Name     string `json:"name"`
	SupportsFormality bool   `json:"supports_formality"` //unused
}


type apiClient struct {
	authKey    string
	baseURL    string
	httpClient *http.Client
}

const (
	baseURLPro  = "https://api.deepl.com/v2"
	baseURLFree = "https://api-free.deepl.com/v2"
)
