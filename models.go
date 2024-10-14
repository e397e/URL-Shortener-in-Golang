package main

import (
	"database/sql"
	"errors"
	"math/rand"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Define custom errors
var (
	ErrURLNotFound = errors.New("URL not found")
	ErrInvalidURL  = errors.New("invalid URL format")
)

var db *sql.DB

// URL represents the URL model
type URL struct {
	ID        int
	ShortCode string
	LongURL   string
	CreatedAt time.Time
}

// InitDB initializes the database connection and creates the table if not exists
func InitDB(filepath string) error {
	var err error
	db, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}

	// Create table if not exists
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS urls (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        short_code TEXT NOT NULL UNIQUE,
        long_url TEXT NOT NULL,
        created_at DATETIME NOT NULL
    );`
	_, err = db.Exec(createTableQuery)
	return err
}

// CloseDB closes the database connection
func CloseDB() {
	db.Close()
}

// CreateURL inserts a new URL into the database
func CreateURL(url URL) error {
	insertQuery := `INSERT INTO urls (short_code, long_url, created_at) VALUES (?, ?, ?)`
	_, err := db.Exec(insertQuery, url.ShortCode, url.LongURL, url.CreatedAt)
	return err
}

// GetURL retrieves a URL from the database by short code
func GetURL(shortCode string) (URL, error) {
	var url URL
	query := `SELECT id, short_code, long_url, created_at FROM urls WHERE short_code = ?`
	row := db.QueryRow(query, shortCode)
	err := row.Scan(&url.ID, &url.ShortCode, &url.LongURL, &url.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return url, ErrURLNotFound
		}
		return url, err
	}
	return url, nil
}

// GenerateShortCode generates a random string of fixed length
func GenerateShortCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6
	seed := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(seed)
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rnd.Intn(len(charset))]
	}
	return string(b)
}

// ValidateURL performs a simple URL validation
func ValidateURL(longURL string) error {
	if !isValidURL(longURL) {
		return ErrInvalidURL
	}
	return nil
}
