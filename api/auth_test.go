package note_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	note "github.com/jun-kagawa/note/api"
)

func TestAuthMiddleware(t *testing.T) {
	t.Run("ValidUserIDCookie", func(t *testing.T) {
		userID := uuid.New()
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotID, err := note.GetUserID(r.Context())
			if err != nil {
				t.Errorf("expected user ID in context, got error: %v", err)
			}
			if gotID != userID {
				t.Errorf("expected user ID %v, got %v", userID, gotID)
			}
			w.WriteHeader(http.StatusOK)
		})

		middleware := note.AuthMiddleware(nextHandler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "user_id", Value: userID.String()})
		rec := httptest.NewRecorder()

		middleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}
	})

	t.Run("MissingUserIDCookie", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("next handler should not be called")
		})

		middleware := note.AuthMiddleware(nextHandler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		middleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("InvalidUserIDCookie", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("next handler should not be called")
		})

		middleware := note.AuthMiddleware(nextHandler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "user_id", Value: "invalid-uuid"})
		rec := httptest.NewRecorder()

		middleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("PostRequestMissingCookie", func(t *testing.T) {
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Error("next handler should not be called")
		})

		middleware := note.AuthMiddleware(nextHandler)

		req := httptest.NewRequest(http.MethodPost, "/notes", nil)
		rec := httptest.NewRecorder()

		middleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}
