package deepl

import (
	"bytes"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func get(endpoint string) (responseBody []byte, err error) {
	httpReq, _ := http.NewRequest("GET", client.baseURL+endpoint, nil)
	httpReq.Header.Set("Authorization", "DeepL-Auth-Key "+client.authKey)

	resp, err := client.httpClient.Do(httpReq)
	if err != nil {
		log.Error().Err(err).Str("endpoint", endpoint).Msg("Failed to send GET request to DeepL API")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Str("response", resp.Status).Str("endpoint", endpoint).Msg("Failed to send GET request to DeepL API, non 200 HTTP response")
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Str("endpoint", endpoint).Msg("Failed to read DeepL API response")
		return nil, err
	}

	return respBody, nil
}

func post(endpoint string, data []byte) (responseBody []byte, err error) {
	httpReq, _ := http.NewRequest("POST", client.baseURL+endpoint, bytes.NewReader(data))
	httpReq.Header.Set("Authorization", "DeepL-Auth-Key "+client.authKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(httpReq)
	if err != nil {
		log.Error().Err(err).Str("endpoint", endpoint).Msg("Failed to send request to POST DeepL API")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Str("response", resp.Status).Str("endpoint", endpoint).Msg("Failed to send POST request to DeepL API, non 200 HTTP response")
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Str("endpoint", endpoint).Msg("Failed to read DeepL API response")
		return nil, err
	}

	return respBody, nil
}
