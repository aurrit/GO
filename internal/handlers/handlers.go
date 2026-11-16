package handlers

import (
    "time"
    "encoding/hex"
    "crypto/rand"
    "html/template"
    "net/http"
    "strconv"
    "strings"
    "go-simple-server/internal/models"
    "go-simple-server/internal/repository"
    "golang.org/x/crypto/bcrypt"
)

var (
    sessionRepo *repository.SessionRepository
    repo         *repository.ExpenseRepository
    categoryRepo *repository.CategoryRepository
    templates    = template.Must(template.ParseGlob("web/templates/*.html"))
    userRepo *repository.UserRepository
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

    expenses, err := repo.GetAll(currentUserID(r), filter)
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

    err := repo.Create(currentUserID(r), models.Expense{
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
    if err := repo.Delete(currentUserID(r), id); err != nil {
        http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditExpenseFormHandler(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.URL.Query().Get("id"))
    expense, err := repo.GetByID(currentUserID(r), id)
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

    err := repo.Update(currentUserID(r), models.Expense{
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
func SetUserRepository(u *repository.UserRepository) {
    userRepo = u
}

func renderAuth(w http.ResponseWriter, name, errMsg string) {
    templates.ExecuteTemplate(w, name, struct{ Error string }{errMsg})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        renderAuth(w, "register", "")
        return
    }
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    r.ParseForm()
    email := strings.TrimSpace(r.FormValue("email"))
    password := r.FormValue("password")

    if email == "" || len(password) < 6 {
        renderAuth(w, "register", "Введите email и пароль не короче 6 символов")
        return
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
        return
    }

    if err := userRepo.Create(email, string(hash)); err != nil {
        renderAuth(w, "register", "Пользователь с таким email уже существует")
        return
    }

    http.Redirect(w, r, "/login", http.StatusSeeOther)
}
func SetSessionRepository(s *repository.SessionRepository) {
    sessionRepo = s
}

func newSessionID() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        renderAuth(w, "login", "")
        return
    }
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    r.ParseForm()
    email := strings.TrimSpace(r.FormValue("email"))
    password := r.FormValue("password")

    user, err := userRepo.GetByEmail(email)
    if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
        renderAuth(w, "login", "Неверный email или пароль")
        return
    }

    sessionID, err := newSessionID()
    if err != nil {
        http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
        return
    }
    expires := time.Now().Add(24 * time.Hour)
    if err := sessionRepo.Create(sessionID, user.ID, expires); err != nil {
        http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "session_id",
        Value:    sessionID,
        Path:     "/",
        Expires:  expires,
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
    })
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    if c, err := r.Cookie("session_id"); err == nil {
        sessionRepo.Delete(c.Value)
    }
    http.SetCookie(w, &http.Cookie{
        Name:   "session_id",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    })
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}