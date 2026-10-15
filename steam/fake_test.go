package steam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"watn3y/steamsalty/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	testSteamID uint64 = 76561199223098853
	testAPIKey         = "TESTAPIKEY0123456789"
	testChatID  int64  = -1001234567890
)

// fakeResponse is either an HTTP response or, if err is set, a transport error.
type fakeResponse struct {
	status int
	body   string
	err    error
}

func (fr fakeResponse) toHTTP(r *http.Request) (*http.Response, error) {
	if fr.err != nil {
		return nil, fr.err
	}
	return &http.Response{
		StatusCode: fr.status,
		Status:     fmt.Sprintf("%d %s", fr.status, http.StatusText(fr.status)),
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(fr.body)),
		Request:    r,
	}, nil
}

// fakeSteam is an http.RoundTripper that stands in for steamcommunity.com,
// the Steam Web API and the Telegram Bot API.
type fakeSteam struct {
	mu        sync.Mutex
	community []fakeResponse // steamcommunity.com responses, consumed one per request
	api       fakeResponse   // returned for every api.steampowered.com request
	apiCalls  int
	messages  []string // text of every Telegram sendMessage call

	// DeepL: returned for POST /v2/translate. /v2/languages always offers
	// source RU and target EN-US.
	translation fakeResponse

	exhausted     chan struct{} // closed once all community responses are consumed
	exhaustedOnce sync.Once

	logs *syncBuffer
}

// newFakeSteam routes all HTTP traffic of the steam package (and go-steamapi) to
// a new fake, sets test config values and captures logs. Everything is restored
// when the test ends.
func newFakeSteam(t *testing.T) *fakeSteam {
	t.Helper()
	f := &fakeSteam{exhausted: make(chan struct{}), logs: &syncBuffer{}}

	origTransport, origConfig, origLogger := http.DefaultTransport, config.BotConfig, log.Logger
	http.DefaultTransport = f
	config.BotConfig.Steam.APIKey = testAPIKey
	config.BotConfig.Telegram.ChatID = testChatID
	config.BotConfig.DeepL.Enabled = false
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: f.logs, NoColor: true})

	t.Cleanup(func() {
		if t.Failed() {
			t.Logf("captured logs:\n%s", f.logs.String())
		}
		http.DefaultTransport, config.BotConfig, log.Logger = origTransport, origConfig, origLogger
	})
	return f
}

func (f *fakeSteam) RoundTrip(r *http.Request) (*http.Response, error) {
	switch r.URL.Host {
	case "steamcommunity.com":
		f.mu.Lock()
		if len(f.community) == 0 {
			f.mu.Unlock()
			f.exhaustedOnce.Do(func() { close(f.exhausted) })
			// The watcher loops forever and has no stop mechanism. Park it here
			// for good so it can't touch state (or log) after its test ended.
			select {}
		}
		next := f.community[0]
		f.community = f.community[1:]
		f.mu.Unlock()
		return next.toHTTP(r)

	case "api.steampowered.com":
		f.mu.Lock()
		f.apiCalls++
		resp := f.api
		f.mu.Unlock()
		return resp.toHTTP(r)

	case "api.telegram.org":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		form, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		f.mu.Lock()
		f.messages = append(f.messages, form.Get("text"))
		f.mu.Unlock()
		return fakeResponse{status: http.StatusOK, body: `{"ok":true,"result":{"message_id":1,"chat":{"id":1}}}`}.toHTTP(r)

	case "api.deepl.com", "api-free.deepl.com":
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v2/languages" && r.URL.Query().Get("type") == "source":
			// Deliberately contains HTML special characters to check escaping.
			return fakeResponse{status: http.StatusOK, body: `[{"language":"RU","name":"Russian & <ru>"}]`}.toHTTP(r)
		case r.Method == http.MethodGet && r.URL.Path == "/v2/languages":
			return fakeResponse{status: http.StatusOK, body: `[{"language":"EN-US","name":"English (American)"}]`}.toHTTP(r)
		case r.Method == http.MethodPost && r.URL.Path == "/v2/translate":
			f.mu.Lock()
			resp := f.translation
			f.mu.Unlock()
			return resp.toHTTP(r)
		}
	}
	return nil, fmt.Errorf("fakeSteam: unexpected request to %s", r.URL)
}

func (f *fakeSteam) sentMessages() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.messages...)
}

func (f *fakeSteam) apiCallCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.apiCalls
}

// syncBuffer is a goroutine-safe log sink.
type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *syncBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// commentsPage builds a successful steamcommunity.com comments response.
func commentsPage(t *testing.T, timeLastPost int64, commentsHTML ...string) fakeResponse {
	t.Helper()
	body, err := json.Marshal(CommentsPage{
		Success:      true,
		TotalCount:   len(commentsHTML),
		CommentsHTML: strings.Join(commentsHTML, "\n"),
		TimeLastPost: timeLastPost,
	})
	if err != nil {
		t.Fatal(err)
	}
	return fakeResponse{status: http.StatusOK, body: string(body)}
}

// commentHTML renders one comment with the structure steamcommunity.com currently
// uses, including the empty timestamp element that precedes the real one.
func commentHTML(id uint64, timestamp int64, author, text string) string {
	return fmt.Sprintf(`<div class="commentthread_comment responsive_body_text   " id="comment_%d">
	<div class="commentthread_comment_content">
		<div class="commentthread_comment_author">
			<div class="commentthread_comment_avatar playerAvatar offline"><a href="https://steamcommunity.com/id/%s"><img src="avatar.jpg"></a></div>
			<div class="author_name_group">
				<div class="flex_row">
					<a class="hoverunderline commentthread_author_link" href="https://steamcommunity.com/id/%s"><bdi>%s</bdi></a>
				</div>
				<div class="commentthread_comment_timestamp">
				</div>
				<div class="commentthread_comment_timestamp" title="some date" data-timestamp="%d">just now</div>
			</div>
		</div>
		<div class="commentthread_comment_text" id="comment_content_%d">
			%s		</div>
	</div>
</div>`, id, author, author, author, timestamp, id, text)
}

// playerSummaries builds a Steam Web API GetPlayerSummaries response for the test profile.
func playerSummaries(personaName string) fakeResponse {
	return fakeResponse{status: http.StatusOK, body: fmt.Sprintf(
		`{"response":{"players":[{"steamid":"%d","personaname":"%s","profileurl":"https://steamcommunity.com/id/owner/"}]}}`,
		testSteamID, personaName)}
}

var (
	// Shape of the production failures: an HTML error page that starts with a
	// UTF-8 BOM (0xEF, logged as 'ï' by the JSON decoder).
	steam502 = fakeResponse{status: http.StatusBadGateway, body: "\xEF\xBB\xBF<html><body>Bad Gateway</body></html>"}
	steam503 = fakeResponse{status: http.StatusServiceUnavailable, body: "<html><body>Service Unavailable</body></html>"}
	// What the Web API returns for a missing or invalid key.
	webAPI403 = fakeResponse{status: http.StatusForbidden, body: `<html><head><title>Forbidden</title></head><body><h1>Forbidden</h1>Access is denied.</body></html>`}
)
