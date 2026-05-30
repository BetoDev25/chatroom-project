package main

import (
	"net/http"
	"strings"

	"github.com/BetoDev25/chatroom-project/internal/cookies"
)

func (cfg *apiConfig) middlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check if the user currently has a session_token (is logged in)

		//paths that don't require authentification
		publicPaths := []string{
			"/login.html",
			"/signup.html",
			"/static/",
		}

		//check if current path is public
		for _, path := range publicPaths {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		sessionToken, err := cookies.Read(r, "session_token")
		if err != nil || sessionToken == "" {
			http.Redirect(w, r, "/login.html", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
