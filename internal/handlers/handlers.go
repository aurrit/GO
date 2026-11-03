package handlers

import (
    "html/template"
    "net/http"
    "strconv"

    "go-simple-server/internal/models"
    "go-simple-server/internal/repository"
)

var (
    repo         *repository.ExpenseRepository
    categoryRepo *repository.CategoryRepository
    templates    = template.Must(template.ParseGlob("web/templates/*.html"))
)

func SetRepositories(e *repository.ExpenseRepository, c *repository.CategoryRepository) {
    repo = e
    categoryRepo = c
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    categoryID, _ := strconv.Atoi(q.Get("category"))
    filter := models.ExpenseFilter{
        CategoryID: categoryID,
        DateFrom:   q.Get("from"),
        DateTo:     q.Get("to"),
    }

    expenses, err := repo.GetAll(filter)
    if err != nil {
        http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
        return
    }
    categories, err := categoryRepo.GetAll()
    if err != nil {
        http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
        return
    }

    data := struct {
        Expenses   []models.Expense
        Categories []models.Category
        Filter     models.ExpenseFilter
    }{expenses, categories, filter}
    templates.ExecuteTemplate(w, "layout", data)
}

func AddExpenseHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    r.ParseForm()
    amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
    categoryID, _ := strconv.Atoi(r.FormValue("category_id"))

    err := repo.Create(models.Expense{
        Amount:      amount,
        Description: r.FormValue("description"),
        Date:        r.FormValue("date"),
        CategoryID:  categoryID,
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
    categories, err := categoryRepo.GetAll()
    if err != nil {
        http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
        return
    }

    data := struct {
        Expense    models.Expense
        Categories []models.Category
    }{expense, categories}
    templates.ExecuteTemplate(w, "edit", data)
}

func UpdateExpenseHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    r.ParseForm()
    id, _ := strconv.Atoi(r.FormValue("id"))
    amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
    categoryID, _ := strconv.Atoi(r.FormValue("category_id"))

    err := repo.Update(models.Expense{
        ID:          id,
        Amount:      amount,
        Description: r.FormValue("description"),
        Date:        r.FormValue("date"),
        CategoryID:  categoryID,
    })
    if err != nil {
        http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}