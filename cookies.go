package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

func (cfg *apiConfig) cookieHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprintf(w, "No session_token cookie found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error reading cookie")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message":      "Login successful",
		"cookie_value": sessionCookie.Value,
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
	http.SetCookie(w, sessionCookie)

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

func generateSessionToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
