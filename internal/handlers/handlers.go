package handlers

import (
    "html/template"
    "net/http"
    "strconv"

    "go-simple-server/internal/models"
    "go-simple-server/internal/repository"
)

var (
    repo      *repository.ExpenseRepository
    templates = template.Must(template.ParseGlob("web/templates/*.html"))
)

func SetRepository(r *repository.ExpenseRepository) {
    repo = r
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    expenses, err := repo.GetAll()
    if err != nil {
        http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
        return
    }

    data := struct {
        Expenses []models.Expense
    }{
        Expenses: expenses,
    }
    templates.ExecuteTemplate(w, "layout", data)
}

func AddExpenseHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    r.ParseForm()
    amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
    description := r.FormValue("description")
    date := r.FormValue("date")

    err := repo.Create(models.Expense{
        Amount:      amount,
        Description: description,
        Date:        date,
    })
    if err != nil {
        http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteExpenseHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    id, _ := strconv.Atoi(r.FormValue("id"))
    if err := repo.Delete(id); err != nil {
        http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditExpenseFormHandler(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.URL.Query().Get("id"))
    expense, err := repo.GetByID(id)
    if err != nil {
        http.Error(w, "Трата не найдена", http.StatusNotFound)
        return
    }

    templates.ExecuteTemplate(w, "edit", expense)
}

func UpdateExpenseHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    r.ParseForm()
    id, _ := strconv.Atoi(r.FormValue("id"))
    amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
    description := r.FormValue("description")
    date := r.FormValue("date")

    err := repo.Update(models.Expense{
        ID:          id,
        Amount:      amount,
        Description: description,
        Date:        date,
    })
    if err != nil {
        http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}