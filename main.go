package main

import (
    "database/sql"
    "fmt"
    "net/http"

    _ "modernc.org/sqlite"

    "go-simple-server/internal/handlers"
    "go-simple-server/internal/repository"
    "go-simple-server/internal/database"
)

func main() {
    db, err := sql.Open("sqlite", "expenses.db?_pragma=foreign_keys(1)")
    if err != nil {
        fmt.Printf("Ошибка открытия базы данных: %v\n", err)
        return
    }
    defer db.Close()

       if err := database.Migrate(db, "migrations"); err != nil {
        fmt.Printf("Ошибка миграций: %v\n", err)
        return
    }

    handlers.SetRepositories(
        repository.NewExpenseRepository(db),
        repository.NewCategoryRepository(db),
    )

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