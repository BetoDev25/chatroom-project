package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/BetoDev25/chatroom-project/internal/auth"
	"github.com/BetoDev25/chatroom-project/internal/cookies"
	"github.com/BetoDev25/chatroom-project/internal/database"
)

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	input := params{}
	err := decoder.Decode(&input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode input")
		return
	}

	user, err := cfg.db.GetUserByUsername(r.Context(), input.Username)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	isValid, err := auth.CheckPasswordHash(input.Password, user.HashedPassword)
	if err != nil || !isValid {
		respondWithError(w, http.StatusUnauthorized, "Incorrect username or password")
		return
	}

	token := auth.GenerateSessionToken()
	session, err := cfg.db.CreateSession(r.Context(), database.CreateSessionParams{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(8 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create session")
		return
	}

	maxAge := int(time.Until(session.ExpiresAt).Seconds())
	if maxAge <= 0 {
		maxAge = 60
	}

	//Secure session cookie
	sessionCookie := &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Path:    "/",
		Expires: session.ExpiresAt,
		MaxAge:  maxAge,
		//MaxAge:   int(time.Until(session.ExpiresAt).Seconds()),
		HttpOnly: true,
		//Domain: temporarily excluded for localhost testing
		Secure:   false,                //temporary for localhost testing
		SameSite: http.SameSiteLaxMode, //temporary for localhost testing
	}
	err = cookies.Write(w, *sessionCookie)
	if err != nil {
		log.Printf("Error setting cookie: %v", err)
		respondWithError(w, http.StatusInternalServerError, "server error")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"username": user.Username,
		"message":  "Login successful",
	})
}
