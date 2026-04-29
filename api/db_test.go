package note_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

func setupDB(t *testing.T) *pgx.Conn {
	t.Helper()
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgresql://note:note@localhost:6432/note_test")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	if err := conn.Ping(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}
	t.Cleanup(func() {
		conn.Close(ctx)
	})
	return conn
}
