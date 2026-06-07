package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteRoom(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	roomNameStr := r.PathValue("roomName")

	room, err := cfg.db.GetRoomByName(r.Context(), roomNameStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "room not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "could not find room")
		}
		return
	}

	if room.OwnerID != userID {
		respondWithError(w, http.StatusForbidden, "User is not the owner of this room")
		return
	}

	err = cfg.db.DeleteRoom(r.Context(), room.RoomID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete room")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Room successfully deleted",
	})
}
