package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/srinivassivaratri/Chirpy/internal/auth"
	"github.com/srinivassivaratri/Chirpy/internal/database"
)

// Defines what a Chirp looks like - it's like a template for tweets
type Chirp struct {
	ID        uuid.UUID `json:"id"`         // Unique identifier for each chirp, like a serial number
	CreatedAt time.Time `json:"created_at"` // When the chirp was first made
	UpdatedAt time.Time `json:"updated_at"` // When the chirp was last changed
	UserID    uuid.UUID `json:"user_id"`    // ID of user who made the chirp
	Body      string    `json:"body"`       // The actual message content
}

// Function that handles creating new chirps when users post them
func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	// Define what data we expect from the user - just the message body
	type parameters struct {
		Body string `json:"body"`
	}

	// Get the user's login token from the request header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		// If no token found, tell user they need to login
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	// Check if the token is valid and get the user's ID
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		// If token invalid, tell user to login again
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	// Set up to read the JSON data sent by user
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		// If can't understand the data, report error
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Check if message follows rules (length, no bad words)
	cleaned, err := validateChirp(params.Body)
	if err != nil {
		// If breaks rules, tell user what's wrong
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	// Try to save the chirp in database
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: userID,
		Body:   cleaned,
	})
	if err != nil {
		// If saving fails, report error
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	// Success! Send back the created chirp details
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID:    chirp.UserID,
		Body:      chirp.Body,
	})
}

// Function to check if a chirp follows the rules
func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	// Check if message is too long (more than 140 characters)
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	// List of words not allowed
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	// Clean the message by replacing bad words
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

// Function that replaces bad words with ****
func getCleanedBody(body string, badWords map[string]struct{}) string {
	// Split message into individual words
	words := strings.Split(body, " ")
	// Check each word
	for i, word := range words {
		// Convert to lowercase to catch variations like "BAD" or "bad"
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			// If it's a bad word, replace with ****
			words[i] = "****"
		}
	}
	// Put all words back together into a message
	cleaned := strings.Join(words, " ")
	return cleaned
}
