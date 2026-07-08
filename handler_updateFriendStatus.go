package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/BetoDev25/chatroom-project/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateFriendStatus(w http.ResponseWriter, r *http.Request) {
	type params struct {
		FriendshipID string `json:"friendship_id"`
		Status       string `json:"status"`
	}

	decoder := json.NewDecoder(r.Body)
	input := params{}
	err := decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode input", err)
		return
	}

	friendshipID, err := uuid.Parse(input.FriendshipID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid friendship ID format", err)
		return
	}

	message := "Friend request accepted"
	if input.Status == "accepted" {
		log.Printf("Updating friendship ID: %s, status: %s", friendshipID, input.Status)
		err = cfg.db.UpdateFriendStatus(r.Context(), database.UpdateFriendStatusParams{
			FriendshipID: friendshipID,
			FriendStatus: input.Status,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				respondWithError(w, http.StatusConflict, "Already friends", err)
				return
			}
			respondWithError(w, http.StatusInternalServerError, "Couldn't add friend", err)
			return
		}
	} else if input.Status == "rejected" {
		message = "Friend request rejected"
		err = cfg.db.DeleteFriendship(r.Context(), friendshipID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't delete friend request", err)
			return
		}
	} else {
		respondWithError(w, http.StatusBadRequest, "Invalid status. Must be 'accepted' or 'rejected'", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": message,
		"status":  input.Status,
	})
}
