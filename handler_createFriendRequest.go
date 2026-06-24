package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/BetoDev25/chatroom-project/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateFriendRequest(w http.ResponseWriter, r *http.Request) {
	type params struct {
		SenderID   uuid.UUID `json:"sender_id"`
		ReceiverID uuid.UUID `json:"receiver_id"`
	}

	decoder := json.NewDecoder(r.Body)
	input := params{}
	err := decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode input")
		return
	}

	request, err := cfg.db.CreateFriendRequest(r.Context(), database.CreateFriendRequestParams{
		SenderID:   input.SenderID,
		ReceiverID: input.ReceiverID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			respondWithError(w, http.StatusConflict, "friend request already exists, or you are already friends with user")
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't send friend request")
		return
	}

	respondWithJSON(w, http.StatusOK, request)
}
