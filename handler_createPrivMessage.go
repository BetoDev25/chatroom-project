package main

import (
	"encoding/json"
	"net/http"

	"github.com/BetoDev25/chatroom-project/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreatePrivMessage(w http.ResponseWriter, r *http.Request) {
	type params struct {
		ConversationID   uuid.UUID `json:"conversation_id"`
		EncryptedContent string    `json:"content"`
	}
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	decoder := json.NewDecoder(r.Body)
	input := params{}
	err := decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode input")
		return
	}

	message, err := cfg.db.CreatePersonalMessage(r.Context(), database.CreatePersonalMessageParams{
		ConversationID:   input.ConversationID,
		UserID:           userID,
		EncryptedContent: input.EncryptedContent,
		MessageType:      "string",
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save message")
		return
	}

	respondWithJSON(w, http.StatusOK, message)
}
