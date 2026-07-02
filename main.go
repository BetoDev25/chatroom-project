package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BetoDev25/chatroom-project/internal/database"
	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development - restrict this in production
	//TO-DO: Figure out what this means
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "guest"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
	}
	client.hub.register <- client

	//read and write pumps in separate goroutines.
	go client.writePump()
	go client.readPump()
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

	//setup Hub
	hub := newHub()
	go hub.run()

	//setup Routes
	mux := http.NewServeMux()
	apiCfg := apiConfig{
		db: dbQueries,
	}
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	// API routes
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)
	mux.HandleFunc("POST /api/cookie", apiCfg.setCookieHandler)
	mux.HandleFunc("GET /api/me", apiCfg.handlerMe)
	mux.HandleFunc("POST /api/logout", apiCfg.handlerLogoutUser)
	mux.HandleFunc("GET /api/users/{username}", apiCfg.handlerGetUserByName)
	mux.HandleFunc("POST /api/rooms/{roomName}", apiCfg.middlewareFunc(apiCfg.handlerCreateRoom))
	mux.HandleFunc("DELETE /api/rooms/{roomName}", apiCfg.middlewareFunc(apiCfg.handlerDeleteRoom))
	mux.HandleFunc("GET /api/rooms/{roomName}", apiCfg.handlerGetRoom)
	mux.HandleFunc("GET /api/rooms", apiCfg.middlewareFunc(apiCfg.handlerGetRooms))
	mux.HandleFunc("POST /api/messages", apiCfg.handlerCreateMessage)
	mux.HandleFunc("GET /api/messages/{roomID}", apiCfg.handlerGetMessages)
	mux.HandleFunc("POST /api/friend-request", apiCfg.handlerCreateFriendRequest)
	mux.HandleFunc("PATCH /api/friend-request", apiCfg.handlerUpdateFriendStatus)
	mux.HandleFunc("GET /api/friend-request/{status}", apiCfg.middlewareFunc(apiCfg.handlerGetRequests))
	mux.HandleFunc("GET /api/conversations/{friendshipID}", apiCfg.middlewareFunc(apiCfg.handlerGetConvo))
	mux.HandleFunc("POST /api/conversations/{friendshipID}", apiCfg.handlerCreateConvo)
	mux.HandleFunc("POST /api/priv-messages", apiCfg.middlewareFunc(apiCfg.handlerCreatePrivMessage))
	mux.HandleFunc("GET /api/priv-messages/{convoID}", apiCfg.middlewareFunc(apiCfg.handlerGetConvoMessages))

	//websocket route
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	mux.HandleFunc("/", apiCfg.middlewareFunc(func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
	}))

	fmt.Println("Server is running on port" + server.Addr)

	fmt.Println("DB_URL:", os.Getenv("DB_URL"))
	// Start server on port specified above
	log.Fatal(server.ListenAndServe())
}
