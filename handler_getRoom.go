package main

import (
	"database/sql"
	"errors"
	"net/http"
)

func (cfg *apiConfig) handlerGetRoom(w http.ResponseWriter, r *http.Request) {
	roomNameStr := r.PathValue("roomName")

	room, err := cfg.db.GetRoomByName(r.Context(), roomNameStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "room not found")
		} else {
			respondWithError(w, http.StatusUnauthorized, "could not find room name")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, room)
}
