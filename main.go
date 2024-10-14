package main

import (
	"log"
	"net/http"
)

func main() {
	// Initialize the database
	err := InitDB("urlshortener.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer CloseDB()

	// Set up routes
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/shorten", ShortenHandler)
	http.HandleFunc("/r/", RedirectHandler)

	// Serve static files if any (optional)
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server
	log.Println("Server started at :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
