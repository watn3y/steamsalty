package steam

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestGetPlayerDetails_Success(t *testing.T) {
	f := newFakeSteam(t)
	f.api = playerSummaries("Owner")

	summary, err := GetPlayerDetails(testSteamID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.SteamID != testSteamID || summary.PersonaName != "Owner" || summary.ProfileURL != "https://steamcommunity.com/id/owner/" {
		t.Fatalf("summary not decoded correctly: %+v", summary)
	}
}

// Before the fix each of these panicked with "index out of range [0] with length 0".
func TestGetPlayerDetails_FailuresReturnError(t *testing.T) {
	tests := []struct {
		name    string
		resp    fakeResponse
		wantErr string
	}{
		{"403 invalid or missing key (production crash)", webAPI403, "Status code 403"},
		{"200 with no players (e.g. deleted account)", fakeResponse{status: http.StatusOK, body: `{"response":{"players":[]}}`}, "no player summary"},
		{"transport error", fakeResponse{err: errors.New("connection reset by peer")}, "connection reset by peer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFakeSteam(t)
			f.api = tt.resp

			_, err := GetPlayerDetails(testSteamID)
			if err == nil {
				t.Fatal("expected an error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error %q does not mention %q", err, tt.wantErr)
			}
			// Transport errors from net/http embed the request URL, which carries key=.
			if strings.Contains(err.Error(), testAPIKey) {
				t.Fatalf("error leaks the Steam API key: %q", err)
			}
		})
	}
}
