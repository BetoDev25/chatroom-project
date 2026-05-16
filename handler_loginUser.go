package main

import (
	"github.com/BetoDev25/chatroom-project/internal/auth"
	"github.com/BetoDev25/chatroom-project/internal/database"
)

func (cfg *apiCfg) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Username string `json:"username"`
		Password string `json:"password`
	}

	decoder := json.NewDecoder(r.Body)
	input := params{}
	err := decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode input")
		return
	}

	user, err := cfg.db.GetUserByUsername(r.Context(), input.Username)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	isValid, err := auth.CheckPasswordHash(input.Password, user.HashedPassword)
	if err != nil || !isValid {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"username": user.Username
		"message": "Login successful",
	})
}
