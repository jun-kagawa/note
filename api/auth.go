package note

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, err := func() (*uuid.UUID, error) {
			cookie, err := r.Cookie("user_id")
			if err != nil {
				return nil, errors.New("no set user_id")
			}
			if v, err := uuid.Parse(cookie.Value); err != nil {
				return nil, err
			} else {
				return &v, nil
			}
		}()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(SetUserID(ctx, *userID)))
	})
}
