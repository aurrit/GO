package repository

import (
    "database/sql"

    "go-simple-server/internal/models"
)

type CategoryRepository struct {
    db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
    return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
    rows, err := r.db.Query("SELECT id, name FROM categories ORDER BY name")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []models.Category
    for rows.Next() {
        var c models.Category
        if err := rows.Scan(&c.ID, &c.Name); err != nil {
            return nil, err
        }
        categories = append(categories, c)
    }
    return categories, nil
}