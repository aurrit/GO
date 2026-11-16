package repository

import (
    "database/sql"
    "time"
)

type SessionRepository struct {
    db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
    return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(id string, userID int, expiresAt time.Time) error {
    _, err := r.db.Exec(
        "INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
        id, userID, expiresAt.UTC().Format(time.RFC3339),
    )
    return err
}

func (r *SessionRepository) GetUserID(id string) (int, error) {
    var userID int
    var expires string
    err := r.db.QueryRow(
        "SELECT user_id, expires_at FROM sessions WHERE id = ?", id,
    ).Scan(&userID, &expires)
    if err != nil {
        return 0, err
    }

    t, err := time.Parse(time.RFC3339, expires)
    if err != nil {
        return 0, err
    }
    if time.Now().After(t) {
        r.Delete(id)
        return 0, sql.ErrNoRows
    }
    return userID, nil
}

func (r *SessionRepository) Delete(id string) error {
    _, err := r.db.Exec("DELETE FROM sessions WHERE id = ?", id)
    return err
}