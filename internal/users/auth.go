package users

import (
	"net/http"
	"strings"
)

func AuthMiddleware(users Users) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			splitAuth := strings.Split(r.Header.Get("Authorization"), "Bearer ")
			// Allow unauthenticated users
			if len(splitAuth) < 2 {
				next.ServeHTTP(w, r)
				return
			}

			user, err := users.Verify(r.Context(), splitAuth[0])
			if err != nil {
				http.Error(w, "Invalid User", http.StatusForbidden)
				return
			}

			// and call the next with our new context
			next.ServeHTTP(w, r.WithContext(users.setCurrent(r.Context(), user)))
		})
	}
}
