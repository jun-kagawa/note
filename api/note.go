package note

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Note struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type NoteOption func(*Note) *Note

func NewNote(userID uuid.UUID, title, body string, opts ...NoteOption) *Note {
	id, _ := uuid.NewV7()
	note := &Note{
		ID:        id,
		UserID:    userID,
		Title:     title,
		Body:      body,
		CreatedAt: time.Now(),
	}
	for _, opt := range opts {
		note = opt(note)
	}
	return note
}

func WithID(ID uuid.UUID) NoteOption {
	return func(note *Note) *Note {
		note.ID = ID
		return note
	}
}

type NoteListItem struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Title     string
	CreatedAt time.Time
}

type NoteRepository struct {
	conn *pgxpool.Pool
}

func NewNoteRepository(conn *pgxpool.Pool) *NoteRepository {
	return &NoteRepository{
		conn: conn,
	}
}

func (r *NoteRepository) Save(ctx context.Context, note *Note) error {
	stmt := "INSERT INTO notes(id, user_id, title, body, created_at) values ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE SET user_id = EXCLUDED.user_id, title = EXCLUDED.title, body = EXCLUDED.body"
	_, err := r.conn.Exec(ctx, stmt, note.ID, note.UserID, note.Title, note.Body, note.CreatedAt)
	return err
}

func (r *NoteRepository) Find(ctx context.Context, ID uuid.UUID) (*Note, error) {
	stmt := "SELECT id, user_id, title, body, created_at FROM notes WHERE id = $1"
	var note Note
	err := r.conn.QueryRow(ctx, stmt, ID).
		Scan(&note.ID, &note.UserID, &note.Title, &note.Body, &note.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &note, err
}

func (r *NoteRepository) Delete(ctx context.Context, ID uuid.UUID) error {
	stmt := "DELETE FROM notes WHERE id = $1"
	_, err := r.conn.Exec(ctx, stmt, ID)
	return err
}

func (r *NoteRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]NoteListItem, error) {
	limit := 10
	stmt := "SELECT id, user_id, title, created_at FROM notes WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2"
	rows, err := r.conn.Query(ctx, stmt, userID, limit)
	if err != nil {
		return nil, err
	}
	items := make([]NoteListItem, 0, limit)
	var id, userId uuid.UUID
	var title string
	var createdAt time.Time
	_, err = pgx.ForEachRow(rows, []any{&id, &userId, &title, &createdAt}, func() error {
		items = append(items, NoteListItem{
			ID:        id,
			UserID:    userId,
			Title:     title,
			CreatedAt: createdAt,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}
