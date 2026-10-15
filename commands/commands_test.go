package commands

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"watn3y/steamsalty/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const testChatID int64 = -1001234567890

// fakeTransport answers the Steam Web API, steamcommunity.com and the Telegram
// Bot API, and records the text of every Telegram message.
type fakeTransport struct {
	mu       sync.Mutex
	messages []string
}

func (f *fakeTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	var body string
	switch r.URL.Host {
	case "api.steampowered.com":
		body = `{"response":{"players":[{"steamid":"76561199223098853","personaname":"Tom & <Jerry>","profileurl":"https://steamcommunity.com/id/tom/"}]}}`
	case "steamcommunity.com":
		body = `{"success":true,"total_count":3,"comments_html":"","timelastpost":1760400000}`
	case "api.telegram.org":
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		form, err := url.ParseQuery(string(raw))
		if err != nil {
			return nil, err
		}
		f.mu.Lock()
		f.messages = append(f.messages, form.Get("text"))
		f.mu.Unlock()
		body = `{"ok":true,"result":{"message_id":2,"chat":{"id":1}}}`
	default:
		return nil, fmt.Errorf("fakeTransport: unexpected request to %s", r.URL)
	}
	return &http.Response{StatusCode: http.StatusOK, Status: "200 OK", Header: http.Header{}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
}

// Persona names are user-controlled. Unescaped, "Tom & <Jerry>" makes Telegram
// reject the whole /info reply.
func TestInfo_EscapesPersonaName(t *testing.T) {
	f := &fakeTransport{}
	origTransport, origConfig, origLogger := http.DefaultTransport, config.BotConfig, log.Logger
	http.DefaultTransport = f
	log.Logger = zerolog.Nop()
	config.BotConfig.Telegram.ChatID = testChatID
	config.BotConfig.Steam.WatchedProfiles = []uint64{76561199223098853}
	t.Cleanup(func() { http.DefaultTransport, config.BotConfig, log.Logger = origTransport, origConfig, origLogger })

	bot := &tgbotapi.BotAPI{Token: "TEST", Client: &http.Client{Transport: f}}
	bot.SetAPIEndpoint(tgbotapi.APIEndpoint)
	update := tgbotapi.Update{Message: &tgbotapi.Message{MessageID: 1, Chat: &tgbotapi.Chat{ID: testChatID}}}

	info(update, bot)

	if len(f.messages) != 1 {
		t.Fatalf("expected exactly 1 reply, got %d: %q", len(f.messages), f.messages)
	}
	want := `<b><a href="https://steamcommunity.com/id/tom/">Tom &amp; &lt;Jerry&gt;</a></b>:`
	if !strings.Contains(f.messages[0], want) {
		t.Fatalf("reply does not contain %q:\n%s", want, f.messages[0])
	}
	if strings.Contains(f.messages[0], "<Jerry>") {
		t.Fatalf("reply contains the unescaped persona name:\n%s", f.messages[0])
	}
}
