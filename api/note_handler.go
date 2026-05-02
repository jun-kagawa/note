package note

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type NoteHandler struct {
	userRepository *UserRepository
	noteRepository *NoteRepository
}

func NewNoteHandler(userRepository *UserRepository, noteRepository *NoteRepository) *NoteHandler {
	return &NoteHandler{
		userRepository: userRepository,
		noteRepository: noteRepository,
	}
}

func (h *NoteHandler) GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := GetUserID(ctx)
	if err != nil {
		httpError(w, http.StatusUnauthorized)
		return
	}
	noteID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpError(w, http.StatusBadRequest)
		return
	}
	user, err := h.userRepository.Find(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpError(w, http.StatusNotFound)
		} else {
			httpError(w, http.StatusInternalServerError)
		}
		return
	}
	note, err := h.noteRepository.Find(ctx, noteID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpError(w, http.StatusNotFound)
		} else {
			httpError(w, http.StatusInternalServerError)
		}
		return
	}
	if note.UserID != user.ID {
		httpError(w, http.StatusForbidden)
		return
	}
	body, err := json.Marshal(note)
	if err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func GetNoteItemHandler(w http.ResponseWriter, r *http.Request) {

}

func UpsertNotHandler(w http.ResponseWriter, r *http.Request) {

}
