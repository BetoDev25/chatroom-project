package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetRooms(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	rooms, err := cfg.db.GetRooms(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusOK, []interface{}{})
		} else {
			respondWithError(w, http.StatusUnauthorized, "could not find room", err)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, rooms)
}
