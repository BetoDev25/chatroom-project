package main

import (
	"encoding/json"
	"net/http"

	"github.com/BetoDev25/chatroom-project/internal/database"
	"github.com/google/uuid"
)

// It's more accurate to say it archives the message, but I want to follow consistent naming convention.
func (cfg *apiConfig) handlerCreateMessage(w http.ResponseWriter, r *http.Request) {
	type params struct {
		RoomID  uuid.UUID `json:"room_id"`
		UserID  uuid.UUID `json:"user_id"`
		Content string    `json:"content"`
	}

	decoder := json.NewDecoder(r.Body)
	input := params{}
	err := decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode input", err)
		return
	}

	message, err := cfg.db.CreateMessage(r.Context(), database.CreateMessageParams{
		RoomID:      input.RoomID,
		UserID:      input.UserID,
		Content:     input.Content,
		MessageType: "string",
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save message", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, message)

}
