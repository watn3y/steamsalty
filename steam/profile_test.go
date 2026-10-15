package steam

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
	"time"
	"watn3y/steamsalty/config"
	deepl "watn3y/steamsalty/deepL"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	t1 int64 = 1760400000
	t2 int64 = 1760400100
)

// runWatcher runs the real watcher against the fake until it has consumed every
// scripted steamcommunity.com response. A panic fails the test.
func runWatcher(t *testing.T, f *fakeSteam) {
	t.Helper()
	bot := &tgbotapi.BotAPI{Token: "TEST", Client: &http.Client{Transport: f}}
	bot.SetAPIEndpoint(tgbotapi.APIEndpoint)

	panicked := make(chan string, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicked <- fmt.Sprintf("%v\n%s", r, debug.Stack())
			}
		}()
		watcher(bot, testSteamID, time.Millisecond)
	}()

	select {
	case <-f.exhausted:
	case p := <-panicked:
		t.Fatalf("watcher panicked: %s", p)
	case <-time.After(30 * time.Second):
		t.Fatal("watcher did not consume the scripted responses in time")
	}
}

// The sequence from the production logs: a single 502 from steamcommunity.com
// while the Web API answers 403. No comment was posted, so nothing may be sent and
// the Web API must not even be called. Before the fix this panicked in GetPlayerDetails.
func TestWatcher_SteamErrorIsNotTreatedAsNewComment(t *testing.T) {
	f := newFakeSteam(t)
	f.api = webAPI403
	baseline := commentsPage(t, t1, commentHTML(1, t1, "old", "old comment"))
	f.community = []fakeResponse{
		baseline,
		steam502,
		steam503,
		{err: errors.New("connection reset by peer")},
		baseline,
	}

	runWatcher(t, f)

	if msgs := f.sentMessages(); len(msgs) != 0 {
		t.Fatalf("expected no notifications, got %d: %q", len(msgs), msgs)
	}
	if n := f.apiCallCount(); n != 0 {
		t.Fatalf("expected no player summary lookups, got %d", n)
	}
	if !strings.Contains(f.logs.String(), "502 Bad Gateway") {
		t.Fatal("expected the 502 to be logged")
	}
}

// A comment posted while Steam was failing must be delivered on the next good poll,
// even when the owner lookup fails (403), using the SteamID instead of the persona name.
// Before the fix the 502 reset the baseline, so this comment was silently dropped.
func TestWatcher_CommentDuringOutageIsNotifiedDespiteOwnerLookupFailure(t *testing.T) {
	f := newFakeSteam(t)
	f.api = webAPI403
	oldComment := commentHTML(1, t1, "old", "old comment")
	f.community = []fakeResponse{
		commentsPage(t, t1, oldComment),
		steam502,
		commentsPage(t, t2, commentHTML(2, t2, "someone", "new comment"), oldComment),
	}

	runWatcher(t, f)

	msgs := f.sentMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected exactly 1 notification, got %d: %q", len(msgs), msgs)
	}
	id := strconv.FormatUint(testSteamID, 10)
	for _, want := range []string{
		"new comment",
		`<a href="https://steamcommunity.com/profiles/` + id + `/">` + id + `</a>'s profile`,
	} {
		if !strings.Contains(msgs[0], want) {
			t.Errorf("notification %q does not contain %q", msgs[0], want)
		}
	}
	if strings.Contains(msgs[0], "old comment") {
		t.Errorf("already-seen comment was notified again: %q", msgs[0])
	}
	if strings.Contains(f.logs.String(), testAPIKey) {
		t.Error("logs leak the Steam API key")
	}
}

// Normal path: a new comment is notified with the owner's persona name from the Web API.
func TestWatcher_NewCommentUsesOwnerPersonaName(t *testing.T) {
	f := newFakeSteam(t)
	f.api = playerSummaries("Owner")
	f.community = []fakeResponse{
		commentsPage(t, t1),
		commentsPage(t, t2, commentHTML(2, t2, "someone", "new comment")),
	}

	runWatcher(t, f)

	msgs := f.sentMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected exactly 1 notification, got %d: %q", len(msgs), msgs)
	}
	want := `<a href="https://steamcommunity.com/id/owner/">Owner</a>'s profile`
	if !strings.Contains(msgs[0], want) || !strings.Contains(msgs[0], "new comment") {
		t.Fatalf("notification %q does not contain %q and the comment text", msgs[0], want)
	}
}

// assertTelegramHTML fails if text has a < or > outside the tags the bot emits
// itself, or an & that doesn't start an entity Telegram accepts. Telegram rejects
// such messages, and raw tags would let commenters inject markup or links.
func assertTelegramHTML(t *testing.T, text string) {
	t.Helper()
	ownTags := regexp.MustCompile(`</?(?:b|blockquote)>|<a href="[^"<>]*">|</a>`)
	if rest := ownTags.ReplaceAllString(text, ""); strings.ContainsAny(rest, "<>") {
		t.Errorf("unescaped < or > in message: %q", text)
	}
	entities := regexp.MustCompile(`&(?:amp|lt|gt|quot|#[0-9]+);`)
	if rest := entities.ReplaceAllString(text, ""); strings.Contains(rest, "&") {
		t.Errorf("unescaped & in message: %q", text)
	}
}

func TestWatcher_EscapesHTMLInNotification(t *testing.T) {
	f := newFakeSteam(t)
	f.api = playerSummaries("Tom & <Jerry>")
	f.translation = fakeResponse{status: http.StatusOK, body: `{"translations":[{"detected_source_language":"RU","text":"you are the best <3 & <b>hi</b>"}]}`}
	config.BotConfig.DeepL.Enabled = true
	config.BotConfig.DeepL.TargetLanguage = "EN-US"
	config.BotConfig.DeepL.FreeTier = true
	deepl.Init()

	// Steam serves comment text and names HTML-encoded; goquery decodes them.
	f.community = []fakeResponse{
		commentsPage(t, t1),
		commentsPage(t, t2, commentHTML(2, t2, "R&amp;D &lt;3", "ты лучший &lt;3 &amp; &lt;b&gt;hi&lt;/b&gt;")),
	}

	runWatcher(t, f)

	msgs := f.sentMessages()
	if len(msgs) != 1 {
		t.Fatalf("expected exactly 1 notification, got %d: %q", len(msgs), msgs)
	}
	assertTelegramHTML(t, msgs[0])
	for _, want := range []string{
		`<a href="https://steamcommunity.com/id/R&amp;D &lt;3">R&amp;D &lt;3</a>`,     // comment author
		`<a href="https://steamcommunity.com/id/owner/">Tom &amp; &lt;Jerry&gt;</a>`,  // profile owner
		`<blockquote>ты лучший &lt;3 &amp; &lt;b&gt;hi&lt;/b&gt;</blockquote>`,        // comment
		`Translated from Russian &amp; &lt;ru&gt;:`,                                   // language name
		`<blockquote>you are the best &lt;3 &amp; &lt;b&gt;hi&lt;/b&gt;</blockquote>`, // translation
	} {
		if !strings.Contains(msgs[0], want) {
			t.Errorf("notification does not contain %q:\n%s", want, msgs[0])
		}
	}
}
