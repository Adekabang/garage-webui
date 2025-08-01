package middleware

import (
	"Adekabang/garage-webui/utils"
	"errors"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	authData := utils.GetEnv("AUTH_USER_PASS", "")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := utils.Session.Get(r, "authenticated")

		if authData == "" {
			next.ServeHTTP(w, r)
			return
		}

		if auth == nil || !auth.(bool) {
			utils.ResponseErrorStatus(w, errors.New("unauthorized"), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
