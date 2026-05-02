package note

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupDB(dsn string) (*pgxpool.Pool, func()) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	return conn, func() {
		conn.Close()
	}
}
