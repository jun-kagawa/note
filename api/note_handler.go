package note

import (
	"context"
	"encoding/json"
	"errors"
	"io"
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

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
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

func (h *NoteHandler) GetNoteItems(w http.ResponseWriter, r *http.Request) {
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

type UpsertNoteDTO struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (h *NoteHandler) UpsertNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := h.currentUser(ctx)
	if err != nil {
		httpError(w, http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var dto UpsertNoteDTO
	if err := json.Unmarshal(body, &dto); err != nil {
		httpError(w, http.StatusBadRequest)
		return
	}
	var note *Note
	if dto.ID == "" {
		note = NewNote(user.ID, dto.Title, dto.Body)
	} else {
		id, err := uuid.Parse(dto.ID)
		if err != nil {
			httpError(w, http.StatusBadRequest)
			return
		}
		note, err = h.noteRepository.Find(ctx, id, user.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				httpError(w, http.StatusNotFound)
			} else {
				httpError(w, http.StatusInternalServerError)
			}
			return
		}
		note.Title = dto.Title
		note.Body = dto.Body
	}

	if err := h.noteRepository.Save(ctx, note); err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := h.currentUser(ctx)
	if err != nil {
		httpError(w, http.StatusUnauthorized)
		return
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpError(w, http.StatusBadRequest)
		return
	}
	if err := h.noteRepository.Delete(ctx, id, user.ID); err != nil {
		httpError(w, http.StatusInternalServerError)
		return
	}
}

func (h *NoteHandler) currentUser(ctx context.Context) (*User, error) {
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	return h.userRepository.Find(ctx, userID)
}
