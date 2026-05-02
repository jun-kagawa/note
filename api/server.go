package note

import (
	"log/slog"
	"net/http"
)

func RunServer() {
	conn, _ := SetupDB("postgresql://note:note@localhost:6432/note")
	userRepository := NewUserRepository(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})
	mux.HandleFunc("POST /users", createUserHandler(userRepository))

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	slog.Info("start server")
	srv.ListenAndServe()
}
