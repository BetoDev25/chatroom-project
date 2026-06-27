package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateConvo(w http.ResponseWriter, r *http.Request) {
	friendshipIDStr := r.PathValue("friendshipID")
	friendshipID, err := uuid.Parse(friendshipIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid friendship ID")
		return
	}

	convo, err := cfg.db.CreateConversation(r.Context(), friendshipID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create convo")
		return
	}

	respondWithJSON(w, http.StatusOK, convo)
}
