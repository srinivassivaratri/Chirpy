package main

import (
	"encoding/json"
	"net/http"

	"github.com/srinivassivaratri/Chirpy/internal/auth"
	"github.com/srinivassivaratri/Chirpy/internal/database"
)

// This function handles HTTP requests to update a user's information
func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	// Define what data we expect from the request - email and password
	type parameters struct {
		Email    string `json:"email"`    // The new email address
		Password string `json:"password"` // The new password
	}

	// Define what data we'll send back in the response - just a User object
	type response struct {
		User // Embed the User type
	}

	// Extract the JWT token from the Authorization header
	// The token proves the user is logged in
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing access token", err)
		return
	}

	// Check if the token is valid and get the user's ID from it
	// This prevents one user from modifying another user's data
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
		return
	}

	// Create a tool to read JSON data from the request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{} // Create an empty container for the data
	// Try to fill the container with data from the request
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Convert the plain password into a secure hashed version
	// This makes it safe to store in the database
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	// Send all the new user data to the database to update the user's record
	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,         // Which user to update
		Email:          params.Email,   // Their new email
		HashedPassword: hashedPassword, // Their new hashed password
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	// Send back the updated user information as JSON
	// We only send non-sensitive data like ID and email, never the password
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        updatedUser.ID,
			Email:     updatedUser.Email,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
		},
	})
}
