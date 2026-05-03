package note

import (
	"context"
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
	noteID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpError(w, http.StatusBadRequest)
		return
	}
	user, err := h.currentUser(ctx)
	if err != nil {
		httpError(w, http.StatusUnauthorized)
		return
	}
	note, err := h.noteRepository.Find(ctx, noteID, user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpError(w, http.StatusNotFound)
		} else {
			httpError(w, http.StatusInternalServerError)
		}
		return
	}
	body, err := json.Marshal(note)
	if err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func (h *NoteHandler) GetNoteItemHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := h.currentUser(ctx)
	if err != nil {
		httpError(w, http.StatusUnauthorized)
		return
	}
	noteItems, err := h.noteRepository.ListByUserID(ctx, user.ID)
	if err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(noteItems)
	if err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func UpsertNotHandler(w http.ResponseWriter, r *http.Request) {

}

func (h *NoteHandler) currentUser(ctx context.Context) (*User, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	return h.userRepository.Find(ctx, userID)
}
