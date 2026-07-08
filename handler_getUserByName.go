package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetUserByName(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		respondWithError(w, http.StatusBadRequest, "Username is required", nil)
		return
	}

	user, err := cfg.db.GetUserByUsername(r.Context(), username)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get user", err)
	}

	respondWithJSON(w, http.StatusOK, user)
}
