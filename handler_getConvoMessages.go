package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/BetoDev25/chatroom-project/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetConvoMessages(w http.ResponseWriter, r *http.Request) {
	convoIDStr := r.PathValue("convoID")

	convoID, err := uuid.Parse(convoIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid convo ID format", err)
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page := 1
	limit := 50

	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	offset := (page - 1) * limit

	messages, err := cfg.db.GetRecentConvo(r.Context(), database.GetRecentConvoParams{
		ConversationID: convoID,
		Limit:          int32(limit),
		Offset:         int32(offset),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithJSON(w, http.StatusOK, []interface{}{})
		} else {
			respondWithError(w, http.StatusUnauthorized, "could not find conversation", err)
		}
		return
	}

	if messages == nil {
		messages = []database.GetRecentConvoRow{}
	}

	respondWithJSON(w, http.StatusOK, messages)
}
