package main

import (
	"database/sql"
	"errors"
	"net/http"
)

func (cfg *apiConfig) handlerGetPublicRooms(w http.ResponseWriter, r *http.Request) {
	roomName := r.PathValue("roomName")

	rooms, err := cfg.db.GetPublicRooms(r.Context(), roomName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusOK, []interface{}{})
		} else {
			respondWithError(w, http.StatusInternalServerError, "failed to fetch rooms", err)
		}
		return
	}

	respondWithJSON(w, http.StatusOK, rooms)
}
