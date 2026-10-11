package main

import (
    "fmt"
    "net/http"

    "go-simple-server/internal/handlers"
)

func main() {
    http.HandleFunc("/", handlers.IndexHandler)
    http.HandleFunc("/add", handlers.AddExpenseHandler)

    fmt.Println("Сервер запущен на http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Printf("Ошибка запуска сервера: %v\n", err)
    }
}