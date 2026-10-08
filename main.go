package main

import (
    "fmt"
    "net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Добро пожаловать! Это мой простой Go-сервер.")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    description := `
Проект: My Simple Go Server
Версия: 1.0.0
Описание: Учебный HTTP-сервер на Go с базовой маршрутизацией.
Технологии: Go, net/http
Автор: КУюжуклу А.В
`
    fmt.Fprintf(w, description)
}

// Новый обработчик для /ping
func pingHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong")
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/about", aboutHandler)
    http.HandleFunc("/ping", pingHandler) // Регистрируем новый маршрут

    fmt.Println("Сервер запущен на http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Printf("Ошибка запуска сервера: %v\n", err)
    }
}

