package note_test

import (
	"context"
	"testing"

	note "github.com/jun-kagawa/note/api"
)

func TestNoteRepository(t *testing.T) {
	conn := setupTestDB(t)
	repo := note.NewNoteRepository(conn)
	userRepo := note.NewUserRepository(conn)
	ctx := context.Background()

	// Helper to create a user for testing notes
	setupUser := func(t *testing.T) *note.User {
		t.Helper()
		user := note.NewUser()
		if err := userRepo.Save(ctx, user); err != nil {
			t.Fatalf("failed to setup user: %v", err)
		}
		return user
	}

	t.Run("save and find note", func(t *testing.T) {
		user := setupUser(t)
		n := note.NewNote(user.ID, "test title", "test body")

		if err := repo.Save(ctx, n); err != nil {
			t.Fatalf("failed to save note: %v", err)
		}

		found, err := repo.Find(ctx, n.ID)
		if err != nil {
			t.Fatalf("failed to find note: %v", err)
		}

		if found.Title != n.Title {
			t.Errorf("got title %q, want %q", found.Title, n.Title)
		}
		if found.Body != n.Body {
			t.Errorf("got body %q, want %q", found.Body, n.Body)
		}
		if found.UserID != n.UserID {
			t.Errorf("got user ID %v, want %v", found.UserID, n.UserID)
		}
	})

	t.Run("update note", func(t *testing.T) {
		user := setupUser(t)
		n := note.NewNote(user.ID, "original title", "original body")
		if err := repo.Save(ctx, n); err != nil {
			t.Fatalf("failed to save note: %v", err)
		}

		n.Title = "updated title"
		n.Body = "updated body"
		if err := repo.Save(ctx, n); err != nil {
			t.Fatalf("failed to update note: %v", err)
		}

		found, err := repo.Find(ctx, n.ID)
		if err != nil {
			t.Fatalf("failed to find note: %v", err)
		}

		if found.Title != "updated title" {
			t.Errorf("got title %q, want %q", found.Title, "updated title")
		}
		if found.Body != "updated body" {
			t.Errorf("got body %q, want %q", found.Body, "updated body")
		}
	})

	t.Run("delete note", func(t *testing.T) {
		user := setupUser(t)
		n := note.NewNote(user.ID, "to be deleted", "...")
		if err := repo.Save(ctx, n); err != nil {
			t.Fatalf("failed to save note: %v", err)
		}

		if err := repo.Delete(ctx, n.ID); err != nil {
			t.Fatalf("failed to delete note: %v", err)
		}

		found, err := repo.Find(ctx, n.ID)
		if err == nil {
			t.Error("expected error when finding deleted note, but got nil")
		}
		if found != nil {
			t.Errorf("expected note to be nil after deletion, but got %v", found)
		}
	})

	t.Run("list notes by user ID", func(t *testing.T) {
		user := setupUser(t)
		notes := []*note.Note{
			note.NewNote(user.ID, "note 1", "body 1"),
			note.NewNote(user.ID, "note 2", "body 2"),
		}

		for _, n := range notes {
			if err := repo.Save(ctx, n); err != nil {
				t.Fatalf("failed to save note: %v", err)
			}
		}

		items, err := repo.ListByUserID(ctx, user.ID)
		if err != nil {
			t.Fatalf("failed to list notes: %v", err)
		}

		if len(items) != len(notes) {
			t.Fatalf("got %d items, want %d", len(items), len(notes))
		}

		// Items should be ordered by CreatedAt DESC
		if items[0].Title != "note 2" {
			t.Errorf("got first item title %q, want %q", items[0].Title, "note 2")
		}
		if items[1].Title != "note 1" {
			t.Errorf("got second item title %q, want %q", items[1].Title, "note 1")
		}
	})
}
