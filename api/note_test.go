package note_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	note "github.com/jun-kagawa/note/api"
)

func v7() uuid.UUID {
	id, _ := uuid.NewV7()
	return id
}

func TestNoteRepository(t *testing.T) {
	conn := openDB()
	repo := note.NewNoteRepository(conn)
	user := note.NewUser()
	note.NewUserRepository(conn).Save(context.Background(), user)

	t.Run("save note", func(t *testing.T) {
		ctx := context.Background()
		note := note.NewNote(user.ID, "test title", "test body")
		err := repo.Save(ctx, note)
		if err != nil {
			t.Errorf("failed save note. err: %v", err)
		}

		foundNote, err := repo.Find(ctx, note.ID)
		if err != nil {
			t.Errorf("failed find note. id: %v", note.ID)
		}
		if note.Title != foundNote.Title || note.Body != foundNote.Body {
			t.Errorf("expect same struct. note: %v, found: %v", note, foundNote)
		}
	})

	t.Run("update note", func(t *testing.T) {
		ctx := context.Background()
		note := note.NewNote(user.ID, "test title", "test body")
		err := repo.Save(ctx, note)
		if err != nil {
			t.Errorf("failed save note. err: %v", err)
		}
		note.Title = "updated title"
		note.Body = "updated body"

		err = repo.Save(ctx, note)
		if err != nil {
			t.Errorf("failed save note. err: %v", err)
		}

		foundNote, err := repo.Find(ctx, note.ID)
		if err != nil {
			t.Errorf("failed find note. id: %v", note.ID)
		}
		if note.Title != foundNote.Title || note.Body != foundNote.Body {
			t.Errorf("expect same struct. note: %v, found: %v", note, foundNote)
		}
	})

	t.Run("delete note", func(t *testing.T) {
		ctx := context.Background()
		note := note.NewNote(user.ID, "test title", "test body")
		err := repo.Save(ctx, note)
		if err != nil {
			t.Errorf("failed save note. err: %v", err)
		}

		err = repo.Delete(ctx, note.ID)
		if err != nil {
			t.Errorf("failed delete note. err: %v", err)
		}

		note, err = repo.Find(ctx, note.ID)
		if note != nil {
			t.Errorf("expected note is nil.")
		}
	})

	t.Run("get note list items", func(t *testing.T) {
		user := note.NewUser()
		note.NewUserRepository(conn).Save(context.Background(), user)
		ctx := context.Background()
		n := note.NewNote(user.ID, "test title", "test body")
		err := repo.Save(ctx, n)
		if err != nil {
			t.Errorf("failed save note. err: %v", err)
		}
		n = note.NewNote(user.ID, "second title", "second body")
		err = repo.Save(ctx, n)
		if err != nil {
			t.Errorf("failed save note. err: %v", err)
		}

		items, err := repo.ListByUserID(ctx, user.ID)
		if err != nil {
			t.Errorf("failed list by user id. err: %v", err)
		}
		if items[0].Title != "second title" || items[1].Title != "test title" {
			t.Errorf("failed. err: %v", items)
		}
	})
}
