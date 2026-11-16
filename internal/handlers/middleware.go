package handlers

import (
    "context"
    "net/http"
)

type contextKey string

const userIDKey contextKey = "userID"

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        c, err := r.Cookie("session_id")
        if err != nil {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        userID, err := sessionRepo.GetUserID(c.Value)
        if err != nil {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        ctx := context.WithValue(r.Context(), userIDKey, userID)
        next(w, r.WithContext(ctx))
    }
}

func currentUserID(r *http.Request) int {
    id, _ := r.Context().Value(userIDKey).(int)
    return id
}