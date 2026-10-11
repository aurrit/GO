package handlers

import (
    "html/template"
    "net/http"
    "strconv"

    "go-simple-server/internal/models"
)

var (
    expenses  []models.Expense
    nextID    = 1
    templates = template.Must(template.ParseGlob("web/templates/*.html"))
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
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

    expenses = append(expenses, models.Expense{
        ID:          nextID,
        Amount:      amount,
        Description: description,
        Date:        date,
    })
    nextID++

    http.Redirect(w, r, "/", http.StatusSeeOther)
}