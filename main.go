package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Chirpy starting on :8080...")
	log.Fatal(server.ListenAndServe())
}
