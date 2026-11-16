package repository

import (
    "database/sql"

    "go-simple-server/internal/models"
)

type ExpenseRepository struct {
    db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
    return &ExpenseRepository{db: db}
}

// 0 превращаем в NULL, чтобы трата могла быть без категории
func nullableCategory(id int) any {
    if id == 0 {
        return nil
    }
    return id
}

func (r *ExpenseRepository) Create(userID int, e models.Expense) error {
    _, err := r.db.Exec(
        "INSERT INTO expenses (amount, description, date, category_id, user_id) VALUES (?, ?, ?, ?, ?)",
        e.Amount, e.Description, e.Date, nullableCategory(e.CategoryID), userID,
    )
    return err
}

func (r *ExpenseRepository) GetAll(userID int, f models.ExpenseFilter) ([]models.Expense, error) {
    query := `SELECT e.id, e.amount, e.description, e.date,
                     COALESCE(e.category_id, 0), COALESCE(c.name, 'Без категории')
              FROM expenses e
              LEFT JOIN categories c ON c.id = e.category_id
              WHERE e.user_id = ?`
    args := []any{userID}

    if f.CategoryID > 0 {
        query += " AND e.category_id = ?"
        args = append(args, f.CategoryID)
    }
    if f.DateFrom != "" {
        query += " AND e.date >= ?"
        args = append(args, f.DateFrom)
    }
    if f.DateTo != "" {
        query += " AND e.date <= ?"
        args = append(args, f.DateTo)
    }
    query += " ORDER BY e.date DESC, e.id DESC"

    rows, err := r.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var expenses []models.Expense
    for rows.Next() {
        var e models.Expense
        if err := rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID, &e.CategoryName); err != nil {
            return nil, err
        }
        expenses = append(expenses, e)
    }
    return expenses, nil
}

func (r *ExpenseRepository) GetByID(userID, id int) (models.Expense, error) {
    var e models.Expense
    err := r.db.QueryRow(
        `SELECT e.id, e.amount, e.description, e.date,
                COALESCE(e.category_id, 0), COALESCE(c.name, 'Без категории')
         FROM expenses e
         LEFT JOIN categories c ON c.id = e.category_id
         WHERE e.id = ? AND e.user_id = ?`, id, userID,
    ).Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID, &e.CategoryName)
    return e, err
}

func (r *ExpenseRepository) Update(userID int, e models.Expense) error {
    _, err := r.db.Exec(
        "UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ? AND user_id = ?",
        e.Amount, e.Description, e.Date, nullableCategory(e.CategoryID), e.ID, userID,
    )
    return err
}

func (r *ExpenseRepository) Delete(userID, id int) error {
    _, err := r.db.Exec("DELETE FROM expenses WHERE id = ? AND user_id = ?", id, userID)
    return err
}