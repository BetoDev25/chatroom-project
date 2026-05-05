package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	db *database.Queries
}

func main() {

	// API routes
	//http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	fmt.Println("Server is running on port" + server.Addr)

	// Start server on port specified above
	log.Fatal(server.ListenAndServe())
}
