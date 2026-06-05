package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/BetoDev25/chatroom-project/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetMessages(w http.ResponseWriter, r *http.Request) {
	roomIDStr := r.PathValue("roomID")

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid room ID format")
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

	messages, err := cfg.db.GetRecentMessages(r.Context(), database.GetRecentMessagesParams{
		RoomID: roomID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "messages not found")
		} else {
			respondWithError(w, http.StatusUnauthorized, "could not find room")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, messages)
}
