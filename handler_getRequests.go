package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	status := r.PathValue("status")

	var requests interface{}
	var err error

	if status == "pending" {
		requests, err = cfg.db.GetPendingRequests(r.Context(), userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithJSON(w, http.StatusOK, []interface{}{})
				return
			}
			respondWithError(w, http.StatusInternalServerError, "could not get friend requests")
			return
		}
	} else if status == "accepted" {
		requests, err = cfg.db.GetFriends(r.Context(), userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithJSON(w, http.StatusOK, []interface{}{})
				return
			}
			respondWithError(w, http.StatusInternalServerError, "could not get friends")
			return
		}
	}

	if requests == nil {
		respondWithJSON(w, http.StatusOK, []interface{}{})
		return
	}
	respondWithJSON(w, http.StatusOK, requests)
}
