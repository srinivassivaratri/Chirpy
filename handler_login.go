package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/srinivassivaratri/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	// Define what data we expect to receive from the login request
	// This creates a template for the JSON data with email, password and optional expiration time
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	// Define what data we'll send back after successful login
	// This includes user info and authentication tokens
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	// Create a tool to read JSON data from the request body
	decoder := json.NewDecoder(r.Body)
	// Make an empty container to store the login details
	params := parameters{}
	// Try to fill the container with data from the request
	err := decoder.Decode(&params)
	if err != nil {
		// If we can't read the data, tell the user something went wrong
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Look up the user in our database using their email
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		// If we can't find the user, don't tell them specifically why - just say login failed
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Check if the password matches what's stored in the database
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		// If password is wrong, give same vague error as before
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Set default token lifetime to 1 hour
	expirationTime := time.Hour
	// If user requested shorter time (between 0 and 3600 seconds), use that instead
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	// Create a new authentication token for this user
	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		// If token creation fails, let the user know
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	// If everything worked, send back the user's info and their new authentication token
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: accessToken,
	})
}
