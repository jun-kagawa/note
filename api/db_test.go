package note_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	note "github.com/jun-kagawa/note/api"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := "postgresql://note:note@localhost:6432/note_test"
	conn, f := note.SetupDB(dsn)
	t.Cleanup(f)
	return conn
}
