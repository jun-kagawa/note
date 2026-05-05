package note

import (
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupServeMux(conn *pgxpool.Pool) *http.ServeMux {
	userRepository := NewUserRepository(conn)
	noteRepository := NewNoteRepository(conn)

	noteHandler := NewNoteHandler(userRepository, noteRepository)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", CreateUserHandler(userRepository))
	mux.Handle("GET /notes/{id}", AuthMiddleware(http.HandlerFunc(noteHandler.GetNote)))
	mux.Handle("GET /notes", AuthMiddleware(http.HandlerFunc(noteHandler.GetNoteItems)))
	mux.Handle("POST /notes", AuthMiddleware(http.HandlerFunc(noteHandler.UpsertNote)))
	mux.Handle("DELETE /notes/{id}", AuthMiddleware(http.HandlerFunc(noteHandler.DeleteNote)))
	return mux
}

func httpError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
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
