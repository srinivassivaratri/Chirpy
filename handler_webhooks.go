package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/srinivassivaratri/Chirpy/internal/auth"
)

// This function is created as part of the apiConfig struct and handles incoming webhook messages from Polka payment service
func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	// Create a template for what data we expect to receive from Polka
	// It's like a form with specific fields that need to be filled out
	type parameters struct {
		// The event field will tell us what happened (like "user upgraded their account")
		Event string `json:"event"`
		// The Data struct contains details about who the event affects
		Data struct {
			// UserID is a unique identifier for the specific user, like a customer number
			UserID uuid.UUID `json:"user_id"`
		}
	}

	// Check if Polka included their secret password (API key) in the request headers
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		// If no password was provided, send back an error saying "you need a password to do this"
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key", err)
		return
	}

	// Compare Polka's password with the one we have stored
	if apiKey != cfg.polkaKey {
		// If passwords don't match, send back an error saying "wrong password"
		respondWithError(w, http.StatusUnauthorized, "API key is invalid", err)
		return
	}

	// Set up a tool to read and understand the JSON data Polka sent us
	decoder := json.NewDecoder(r.Body)
	// Create an empty parameters struct to store the data we read
	params := parameters{}
	// Try to fill our parameters struct with Polka's data
	err = decoder.Decode(&params)
	if err != nil {
		// If we can't understand Polka's data, send back an error
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// We only care about "user.upgraded" events - ignore everything else
	if params.Event != "user.upgraded" {
		// Send back a "received but ignored" response
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Try to update the user's account in our database to premium status
	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		// If we can't find a user with that ID, tell Polka "this user doesn't exist"
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		// If any other database error happens, tell Polka "something went wrong on our end"
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	// If everything worked perfectly, send back a simple "ok, done" response
	w.WriteHeader(http.StatusNoContent)
}
