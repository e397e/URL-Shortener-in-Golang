package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))

// HomeHandler renders the homepage with the URL input form
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	templates.ExecuteTemplate(w, "index.html", nil)
}

// ShortenHandler processes the form submission to shorten URLs
func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL := r.FormValue("long_url")
	if err := ValidateURL(longURL); err != nil {
		if err == ErrInvalidURL {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a new shortened URL
	shortCode := GenerateShortCode()

	// Save to the database
	url := URL{
		ShortCode: shortCode,
		LongURL:   longURL,
		CreatedAt: time.Now(),
	}

	err := CreateURL(url)
	if err != nil {
		log.Println("Error saving URL:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare the short URL
	shortURL := "http://" + r.Host + "/r/" + shortCode

	// Render success page
	templates.ExecuteTemplate(w, "success.html", shortURL)
}

// RedirectHandler handles redirection from short URL to the original long URL
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the short code from the URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	shortCode := pathParts[2]

	// Retrieve the long URL from the database
	url, err := GetURL(shortCode)
	if err != nil {
		if err == ErrURLNotFound {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Redirect to the long URL
	http.Redirect(w, r, url.LongURL, http.StatusMovedPermanently)
}

// isValidURL is a simple URL validation function
func isValidURL(url string) bool {
	// Basic check; can be enhanced
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
