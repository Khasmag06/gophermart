package middlewares

import (
	"context"
	"github.com/Khasmag06/gophermart/internal/helpers"
	"net/http"
	"strconv"
)

type userIDKeyType string

const UserIDKey userIDKeyType = "userID"

func CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err == http.ErrNoCookie {
			http.Error(w, "you are unauthorized", http.StatusUnauthorized)
			return
		}
		id, err := helpers.ExtractIDFromToken(cookie.Value)
		if err != nil {
			http.Error(w, "you are unauthorized", http.StatusUnauthorized)
			return
		}
		userID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
