package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetConvo(w http.ResponseWriter, r *http.Request) {
	friendshipIDStr := r.PathValue("friendshipID")
	friendshipID, err := uuid.Parse(friendshipIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid friend ID", err)
		return
	}

	convo, err := cfg.db.GetConvoByFriendshipID(r.Context(), friendshipID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "conversation not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "could not get convo", err)
		return
	}

	respondWithJSON(w, http.StatusOK, convo)
}
