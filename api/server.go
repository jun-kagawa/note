package note

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupServeMux(conn *pgxpool.Pool) *http.ServeMux {
	userRepository := NewUserRepository(conn)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", CreateUserHandler(userRepository))
	return mux
}

func RunServer() {
	conn, _ := SetupDB("postgresql://note:note@localhost:6432/note")

	mux := SetupServeMux(conn)
	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	slog.Info("start server")
	srv.ListenAndServe()
}
