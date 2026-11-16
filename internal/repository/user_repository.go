package repository

import (
    "database/sql"

    "go-simple-server/internal/models"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, passwordHash string) error {
    _, err := r.db.Exec(
        "INSERT INTO users (email, password_hash) VALUES (?, ?)",
        email, passwordHash,
    )
    return err
}

func (r *UserRepository) GetByEmail(email string) (models.User, error) {
    var u models.User
    err := r.db.QueryRow(
        "SELECT id, email, password_hash FROM users WHERE email = ?", email,
    ).Scan(&u.ID, &u.Email, &u.PasswordHash)
    return u, err
}