package main

import (
	"net/http"

	"github.com/BetoDev25/chatroom-project/internal/cookies"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerMe(w http.ResponseWriter, r *http.Request) {
	type UserInfo struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
	}

	sessionToken, err := cookies.Read(r, "session_token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	user, err := cfg.db.GetUserByCookie(r.Context(), sessionToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired session")
		return
	}

	respondWithJSON(w, http.StatusOK, UserInfo{
		ID:       user.ID,
		Username: user.Username,
	})
}
