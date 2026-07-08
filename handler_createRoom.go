package main

import (
	"net/http"
	"strings"

	"github.com/BetoDev25/chatroom-project/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateRoom(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	roomNameStr := r.PathValue("roomName")

	if valid, err := ValidateRoomName(roomNameStr); !valid {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	room, err := cfg.db.CreateRoom(r.Context(), database.CreateRoomParams{
		OwnerID:  userID,
		RoomName: roomNameStr,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			respondWithError(w, http.StatusConflict, "room already exists", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "could not create room", err)
		return
	}

	respondWithJSON(w, http.StatusOK, room)
}
