package note_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	note "github.com/jun-kagawa/note/api"
)

func TestGetNoteHandler(t *testing.T) {
	conn := setupTestDB(t)
	userRepo := note.NewUserRepository(conn)
	noteRepo := note.NewNoteRepository(conn)
	handler := note.NewNoteHandler(userRepo, noteRepo)

	ctx := context.Background()
	user := note.NewUser()
	if err := userRepo.Save(ctx, user); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	n := note.NewNote(user.ID, "Test Note", "Test Body")
	if err := noteRepo.Save(ctx, n); err != nil {
		t.Fatalf("failed to save note: %v", err)
	}

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes/"+n.ID.String(), nil)
		// Go 1.22+ PathValue support
		req.SetPathValue("id", n.ID.String())

		// Set user ID in context as GetNoteHandler expects it from middleware
		req = req.WithContext(note.SetUserID(req.Context(), user.ID))

		rec := httptest.NewRecorder()
		handler.GetNoteHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		var got note.Note
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if got.ID != n.ID {
			t.Errorf("got note ID %v, want %v", got.ID, n.ID)
		}
		if got.Title != n.Title {
			t.Errorf("got title %q, want %q", got.Title, n.Title)
		}
	})

	t.Run("invalid note id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes/invalid-uuid", nil)
		req.SetPathValue("id", "invalid-uuid")
		req = req.WithContext(note.SetUserID(req.Context(), user.ID))

		rec := httptest.NewRecorder()
		handler.GetNoteHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("note not found", func(t *testing.T) {
		unknownID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/notes/"+unknownID.String(), nil)
		req.SetPathValue("id", unknownID.String())
		req = req.WithContext(note.SetUserID(req.Context(), user.ID))

		rec := httptest.NewRecorder()
		handler.GetNoteHandler(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("unauthorized access (other user's note)", func(t *testing.T) {
		otherUser := note.NewUser()
		if err := userRepo.Save(ctx, otherUser); err != nil {
			t.Fatalf("failed to save other user: %v", err)
		}

		otherNote := note.NewNote(otherUser.ID, "Other User Note", "Body")
		if err := noteRepo.Save(ctx, otherNote); err != nil {
			t.Fatalf("failed to save other note: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/notes/"+otherNote.ID.String(), nil)
		req.SetPathValue("id", otherNote.ID.String())
		// Authenticated as 'user', but trying to access 'otherNote'
		req = req.WithContext(note.SetUserID(req.Context(), user.ID))

		rec := httptest.NewRecorder()
		handler.GetNoteHandler(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("missing user id in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes/"+n.ID.String(), nil)
		req.SetPathValue("id", n.ID.String())
		// No user ID set in context

		rec := httptest.NewRecorder()
		handler.GetNoteHandler(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("routing and middleware integration", func(t *testing.T) {
		mux := note.SetupServeMux(conn)
		req := httptest.NewRequest(http.MethodGet, "/notes/"+n.ID.String(), nil)
		req.AddCookie(&http.Cookie{Name: "user_id", Value: user.ID.String()})

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		var got note.Note
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.ID != n.ID {
			t.Errorf("got note ID %v, want %v", got.ID, n.ID)
		}
	})
}
