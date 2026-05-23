package middleware

import (
	"context"
	"database/sql"
	"net/http"

	"bulletin/internal/models"
)

type contextKey string
const UserIDKey contextKey = "userID"

func Auth(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err == nil {
				userID, err := models.GetUserIDBySession(db, cookie.Value)
				if err == nil {
					ctx := context.WithValue(r.Context(), UserIDKey, userID)
					r = r.WithContext(ctx)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetUserID(r) == 0 {
			http.Redirect(w, r, "/auth/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetUserID(r *http.Request) int {
	id, _ := r.Context().Value(UserIDKey).(int)
	return id
}
