package main

import (
	"net/http"

	"github.com/BetoDev25/chatroom-project/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetConvo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	friendIDStr := r.PathValue("friendshipID")
	friendID, err := uuid.Parse(friendIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid friend ID")
		return
	}

	convo, err := cfg.db.GetConvoBetweenUsers(r.Context(), database.GetConvoBetweenUsersParams{
		SenderID:   userID,
		ReceiverID: friendID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not get convo")
		return
	}

	respondWithJSON(w, http.StatusOK, convo)
}
