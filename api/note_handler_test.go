package note_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	note "github.com/jun-kagawa/note/api"
)

func TestGetNoteHandler(t *testing.T) {
	conn := setupTestDB(t)
	userRepo := note.NewUserRepository(conn)
	noteRepo := note.NewNoteRepository(conn)
	mux := note.SetupServeMux(conn)

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
		if got.Title != n.Title {
			t.Errorf("got title %q, want %q", got.Title, n.Title)
		}
	})

	t.Run("invalid note id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes/invalid-uuid", nil)
		req.AddCookie(&http.Cookie{Name: "user_id", Value: user.ID.String()})

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("note not found", func(t *testing.T) {
		unknownID := uuid.New()
		req := httptest.NewRequest(http.MethodGet, "/notes/"+unknownID.String(), nil)
		req.AddCookie(&http.Cookie{Name: "user_id", Value: user.ID.String()})

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

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
		// Authenticated as 'user', but trying to access 'otherNote'
		req.AddCookie(&http.Cookie{Name: "user_id", Value: user.ID.String()})

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})

	t.Run("missing user id in context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes/"+n.ID.String(), nil)
		// No user_id cookie

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestGetNoteItemHandler(t *testing.T) {
	conn := setupTestDB(t)
	userRepo := note.NewUserRepository(conn)
	noteRepo := note.NewNoteRepository(conn)
	mux := note.SetupServeMux(conn)

	ctx := context.Background()
	user := note.NewUser()
	if err := userRepo.Save(ctx, user); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	notes := []note.Note{
		*note.NewNote(user.ID, "Note 1", "Body 1"),
		*note.NewNote(user.ID, "Note 2", "Body 2"),
	}
	for i := range notes {
		if err := noteRepo.Save(ctx, &notes[i]); err != nil {
			t.Fatalf("failed to save note %d: %v", i, err)
		}
	}

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes", nil)
		req.AddCookie(&http.Cookie{Name: "user_id", Value: user.ID.String()})

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		var got []note.Note
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(got) != len(notes) {
			t.Errorf("got %d notes, want %d", len(got), len(notes))
		}
	})

	t.Run("empty list", func(t *testing.T) {
		userWithoutNotes := note.NewUser()
		if err := userRepo.Save(ctx, userWithoutNotes); err != nil {
			t.Fatalf("failed to save user: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/notes", nil)
		req.AddCookie(&http.Cookie{Name: "user_id", Value: userWithoutNotes.ID.String()})

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		var got []note.Note
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(got) != 0 {
			t.Errorf("got %d notes, want 0", len(got))
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes", nil)
		// No user_id cookie

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})
}

func TestUpsertNoteHandler(t *testing.T) {
	conn := setupTestDB(t)
	userRepo := note.NewUserRepository(conn)
	noteRepo := note.NewNoteRepository(conn)
	mux := note.SetupServeMux(conn)

	ctx := context.Background()
	user := note.NewUser()
	if err := userRepo.Save(ctx, user); err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	t.Run("success create", func(t *testing.T) {
		dto := note.UpsertNoteDTO{
			Title: "New Note",
			Body:  "New Body",
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
		req.AddCookie(&http.Cookie{Name: "user_id", Value: user.ID.String()})

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		// Verify note was created
		notes, err := noteRepo.ListByUserID(ctx, user.ID)
		if err != nil {
			t.Fatalf("failed to list notes: %v", err)
		}
		if len(notes) != 1 {
			t.Errorf("got %d notes, want 1", len(notes))
		}
		if notes[0].Title != dto.Title {
			t.Errorf("got title %q, want %q", notes[0].Title, dto.Title)
		}
	})

	t.Run("success update", func(t *testing.T) {
		n := note.NewNote(user.ID, "Old Title", "Old Body")
		if err := noteRepo.Save(ctx, n); err != nil {
			t.Fatalf("failed to save note: %v", err)
		}

		dto := note.UpsertNoteDTO{
			ID:    n.ID.String(),
			Title: "Updated Title",
			Body:  "Updated Body",
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
		req.AddCookie(&http.Cookie{Name: "user_id", Value: user.ID.String()})

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		// Verify note was updated
		updated, err := noteRepo.Find(ctx, n.ID, user.ID)
		if err != nil {
			t.Fatalf("failed to find note: %v", err)
		}
		if updated.Title != dto.Title {
			t.Errorf("got title %q, want %q", updated.Title, dto.Title)
		}
		if updated.Body != dto.Body {
			t.Errorf("got body %q, want %q", updated.Body, dto.Body)
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/notes", nil)
		// No user_id cookie

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", rec.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader("{invalid}"))
		req.AddCookie(&http.Cookie{Name: "user_id", Value: user.ID.String()})

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})

	t.Run("update non-existent note", func(t *testing.T) {
		dto := note.UpsertNoteDTO{
			ID:    uuid.New().String(),
			Title: "Title",
			Body:  "Body",
		}
		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewReader(body))
		req.AddCookie(&http.Cookie{Name: "user_id", Value: user.ID.String()})

		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}
	})
}
