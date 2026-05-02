package note_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	note "github.com/jun-kagawa/note/api"
)

func TestUserRepository(t *testing.T) {
	conn := setupTestDB(t)
	repo := note.NewUserRepository(conn)
	ctx := context.Background()

	t.Run("save and find user", func(t *testing.T) {
		user := note.NewUser()
		if err := repo.Save(ctx, user); err != nil {
			t.Fatalf("failed to save user: %v", err)
		}

		found, err := repo.Find(ctx, user.ID)
		if err != nil {
			t.Fatalf("failed to find user: %v", err)
		}

		if found.ID != user.ID {
			t.Errorf("got user ID %v, want %v", found.ID, user.ID)
		}
	})

	t.Run("find non-existent user", func(t *testing.T) {
		id := uuid.New()
		user, err := repo.Find(ctx, id)
		if err == nil {
			t.Error("expected error when finding non-existent user, but got nil")
		}
		if user != nil {
			t.Errorf("expected user to be nil, but got %v", user)
		}
	})
}
