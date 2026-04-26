package note_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	note "github.com/jun-kagawa/note/api"
)

func openDB() *pgx.Conn {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgresql://note:note@localhost:6432/note_test")
	if err != nil {
		panic(err)
	}
	if err := conn.Ping(ctx); err != nil {
		panic(err)
	}
	return conn
}

func TestUserRepository(t *testing.T) {
	conn := openDB()
	repo := note.NewUserRepository(conn)

	t.Run("save user", func(t *testing.T) {
		ctx := context.Background()
		user := note.NewUser()
		err := repo.Save(ctx, user)
		if err != nil {
			t.Error("failed save user")
		}
		foundUser, err := repo.Find(ctx, user.ID)
		if user.ID != foundUser.ID {
			t.Errorf("failed find user by id. id: %v", user.ID)
		}
	})

	t.Run("not found user", func(t *testing.T) {
		ctx := context.Background()
		user, err := repo.Find(ctx, uuid.New())
		if err == nil {
			t.Errorf("expected not nil")
		}
		if user != nil {
			t.Errorf("user isn't nil")
		}
	})
}
