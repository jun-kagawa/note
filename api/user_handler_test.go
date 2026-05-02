package note_test

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	note "github.com/jun-kagawa/note/api"
)

func TestCreateUserHandler(t *testing.T) {
	conn := setupTestDB(t)
	mux := note.SetupServeMux(conn)

	req := httptest.NewRequest(http.MethodPost, "/users", nil)

	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	cookies := res.Cookies()
	i := slices.IndexFunc(cookies, func(cookie *http.Cookie) bool {
		if cookie.Name == "user_id" {
			return true
		} else {
			return false
		}
	})
	if i < 0 {
		t.Errorf("expected index above 0, got %d", i)
	}
	if cookies[i].Value == "" {
		t.Error("expected user_id cookie to be set")
	}
}
