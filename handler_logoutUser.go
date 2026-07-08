package main

import (
	"net/http"

	"github.com/BetoDev25/chatroom-project/internal/cookies"
)

func (cfg *apiConfig) handlerLogoutUser(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := cookies.Read(r, "session_token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no active session", err)
		return
	}

	err = cfg.db.DeleteSessionByToken(r.Context(), sessionToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete session token", err)
		return
	}

	nullCookie := cookies.DeleteCookie("session_token")
	err = cookies.Write(w, *nullCookie)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete cookie", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Session deleted",
	})
}
