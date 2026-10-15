package steam

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestGetComments_Success(t *testing.T) {
	f := newFakeSteam(t)
	f.community = []fakeResponse{commentsPage(t, 1760400000, commentHTML(1, 1760400000, "someone", "hi"))}

	page, err := GetComments(testSteamID, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !page.Success || page.TimeLastPost != 1760400000 || page.TotalCount != 1 || page.CommentsHTML == "" {
		t.Fatalf("page not decoded correctly: %+v", page)
	}
}

// Every failed request must come back as an error. Before the fix these all
// returned a zero CommentsPage (TimeLastPost 0), which the watcher read as
// "the newest comment changed".
func TestGetComments_FailuresReturnError(t *testing.T) {
	tests := []struct {
		name    string
		resp    fakeResponse
		wantErr string
	}{
		{"502 with BOM HTML body", steam502, "502"},
		{"503 HTML body", steam503, "503"},
		{"200 with non-JSON body", fakeResponse{status: http.StatusOK, body: "<html></html>"}, "parse comments as JSON"},
		{"200 with success=false", fakeResponse{status: http.StatusOK, body: `{"success":false}`}, "success=false"},
		// Used to panic: resp is nil when http.Get fails, and resp.StatusCode was read anyway.
		{"transport error", fakeResponse{err: errors.New("connection reset by peer")}, "connection reset by peer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFakeSteam(t)
			f.community = []fakeResponse{tt.resp}

			page, err := GetComments(testSteamID, 0, 10)
			if err == nil {
				t.Fatalf("expected an error, got page %+v", page)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error %q does not mention %q", err, tt.wantErr)
			}
			if page != (CommentsPage{}) {
				t.Fatalf("expected zero page alongside the error, got %+v", page)
			}
		})
	}
}

// Steam's markup has two .commentthread_comment_timestamp elements per comment and
// only the second carries data-timestamp. Reading the first made every comment fail
// to parse, so no notification was ever sent.
func TestParseComments_CurrentSteamMarkup(t *testing.T) {
	page := CommentsPage{CommentsHTML: strings.Join([]string{
		commentHTML(2, 1760400100, "newer", "second comment"), // Steam lists newest first
		commentHTML(1, 1760400000, "older", "first comment"),
	}, "\n")}

	got := parseComments(page)

	want := []Comment{
		{ID: 1, Timestamp: 1760400000, Author: "older", AuthorProfileURL: "https://steamcommunity.com/id/older", Text: "first comment"},
		{ID: 2, Timestamp: 1760400100, Author: "newer", AuthorProfileURL: "https://steamcommunity.com/id/newer", Text: "second comment"},
	}
	if len(got) != len(want) {
		t.Fatalf("expected %d comments, got %d: %+v", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("comment %d:\n got  %+v\n want %+v", i, got[i], want[i])
		}
	}
}
