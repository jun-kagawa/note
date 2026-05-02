package note

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func clearCookie(w http.ResponseWriter, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

const USER_ID = "UserID"

func SetUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, USER_ID, userID)
}
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	id := ctx.Value(USER_ID)
	if userID, ok := id.(uuid.UUID); ok {
		return userID, nil
	} else {
		return uuid.UUID{}, errors.New("failed cast any to uuid.UUID")
	}
}

func CreateUserHandler(userRepository *UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie("user_id")
		if err == http.ErrNoCookie {
			user := NewUser()
			if err := userRepository.Save(ctx, user); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			cookie := &http.Cookie{
				Name:     "user_id",
				Value:    user.ID.String(),
				Secure:   true,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		userID, err := uuid.Parse(cookie.Value)
		if err != nil {
			clearCookie(w, "user_id")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := userRepository.Find(r.Context(), userID)
		if err != nil {
			clearCookie(w, "user_id")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		cookie = &http.Cookie{
			Name:     "user_id",
			Value:    user.ID.String(),
			Secure:   true,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}
}
