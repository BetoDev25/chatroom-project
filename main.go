package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BetoDev25/chatroom-project/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorMsg struct {
		Error string `json:"error"`
	}
	resp := errorMsg{
		Error: msg,
	}
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		return
	} else {
		w.WriteHeader(code)
	}
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(payload)
}

type apiConfig struct {
	db *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()
	apiCfg := apiConfig{
		db: dbQueries,
	}
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	// API routes
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	fmt.Println("Server is running on port" + server.Addr)

	// Start server on port specified above
	log.Fatal(server.ListenAndServe())
}
