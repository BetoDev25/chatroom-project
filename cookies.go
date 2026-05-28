package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/BetoDev25/chatroom-project/internal/cookies"
)

func (cfg *apiConfig) cookieHandler(w http.ResponseWriter, r *http.Request) {
	value, err := cookies.Read(r, "session_token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			respondWithError(w, http.StatusBadRequest, "cookie not found")
		case errors.Is(err, cookies.ErrInvalidValue):
			respondWithError(w, http.StatusBadRequest, "invalid cookie")
		default:
			log.Println(err)
			respondWithError(w, http.StatusInternalServerError, "Error reading cookie")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message":      "Login successful",
		"cookie_value": value,
	})
}

func (cfg *apiConfig) setCookieHandler(w http.ResponseWriter, r *http.Request) {
	//Secure session cookie
	sessionCookie := &http.Cookie{
		Name:     "session_token",
		Value:    generateSessionToken(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   86400,
		HttpOnly: true,
		//Domain: temporarily excluded for localhost testing
		Secure:   false,                //temporary for localhost testing
		SameSite: http.SameSiteLaxMode, //temporary for localhost testing
	}
	err := cookies.Write(w, *sessionCookie)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "server error")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message":      "Cookie set successfully",
		"cookie_value": sessionCookie.Value,
	})
}

func (cfg *apiConfig) deleteCookieHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie := &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	}
	http.SetCookie(w, sessionCookie)

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Cookie deleted",
	})
}
