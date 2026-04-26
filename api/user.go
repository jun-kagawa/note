package note

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

func NewUser() *User {
	id, _ := uuid.NewV7()
	return &User{
		ID:        id,
		CreatedAt: time.Now(),
	}
}

type UserRepository struct {
	conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{
		conn: conn,
	}
}

func (r *UserRepository) Save(ctx context.Context, user *User) error {
	_, err := r.conn.Exec(ctx, "INSERT INTO users(id, created_at) VALUES ($1, $2)", user.ID, user.CreatedAt)
	return err
}

func (r *UserRepository) Find(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := r.conn.QueryRow(ctx, "SELECT id, created_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, err
}
