package main

import (
    "database/sql"
    "fmt"
    "net/http"
    "os"

    _ "modernc.org/sqlite"

    "go-simple-server/internal/handlers"
    "go-simple-server/internal/repository"
)

func main() {
    db, err := sql.Open("sqlite", "expenses.db")
    if err != nil {
        fmt.Printf("Ошибка открытия базы данных: %v\n", err)
        return
    }
    defer db.Close()

    migration, err := os.ReadFile("migrations/001_create_expenses.sql")
    if err != nil {
        fmt.Printf("Ошибка чтения миграции: %v\n", err)
        return
    }
    if _, err := db.Exec(string(migration)); err != nil {
        fmt.Printf("Ошибка применения миграции: %v\n", err)
        return
    }

    repo := repository.NewExpenseRepository(db)
    handlers.SetRepository(repo)

    http.HandleFunc("/", handlers.IndexHandler)
    http.HandleFunc("/add", handlers.AddExpenseHandler)
    http.HandleFunc("/delete", handlers.DeleteExpenseHandler)
    http.HandleFunc("/edit", handlers.EditExpenseFormHandler)
    http.HandleFunc("/update", handlers.UpdateExpenseHandler)

    fmt.Println("Сервер запущен на http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Printf("Ошибка запуска сервера: %v\n", err)
    }
}