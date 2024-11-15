package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/srinivassivaratri/Chirpy/internal/auth"
	"github.com/srinivassivaratri/Chirpy/internal/database"
)

// This function handles what happens when someone tries to log in to our system
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
    // This defines what information we expect users to send us when logging in
    // Like a form with email and password fields
    type parameters struct {
        Password string `json:"password"` // The secret word/phrase they use to prove who they are
        Email    string `json:"email"`    // Their email address that identifies them
    }

    // This defines what information we'll send back after they log in successfully
    // Like a welcome package containing their ID card and special passes
    type response struct {
        User                              // All their basic account information
        Token        string `json:"token"` // A temporary pass that lets them do things (expires soon)
        RefreshToken string `json:"refresh_token"` // A special ticket they can use to get new temporary passes
    }

    // Create a tool that can read and understand the login information they sent us
    decoder := json.NewDecoder(r.Body)
    // Make an empty container to hold their login information
    params := parameters{}
    // Try to fill our container with their login information
    err := decoder.Decode(&params)
    if err != nil {
        // If we can't understand what they sent us, tell them there's a problem
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
        return
    }

    // Look in our records for someone with this email address
    user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
    if err != nil {
        // If we can't find anyone with this email, tell them login failed
        // We're intentionally vague to prevent hackers from knowing which part was wrong
        respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
        return
    }

    // Check if the password they gave matches what we have stored
    // Like checking if a key fits the lock
    err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
    if err != nil {
        // If the password doesn't match, tell them login failed
        respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
        return
    }

    // Create a special temporary pass (JWT token) that expires in 1 hour
    // Like a visitor badge that proves they logged in properly
    accessToken, err := auth.MakeJWT(
        user.ID,        // Write their ID on the badge
        cfg.jwtSecret,  // Use our special stamp to make it official
        time.Hour,      // Make it expire in 1 hour
    )
    if err != nil {
        // If we can't create their temporary pass, tell them something went wrong
        respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
        return
    }

    // Create a special long-term pass (refresh token) they can use later
    // Like a membership card that lets them get new visitor badges
    refreshToken, err := auth.MakeRefreshToken()
    if err != nil {
        // If we can't create their long-term pass, tell them something went wrong
        respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
        return
    }

    // Save their long-term pass in our records so we remember it's valid
    // Like keeping a copy of their membership card number in our system
    _, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
        Token:  refreshToken,  // The membership card number
        UserID: user.ID,      // Who it belongs to
    })
    if err != nil {
        // If we can't save their long-term pass, tell them something went wrong
        respondWithError(w, http.StatusInternalServerError, "Couldn't store refresh token", err)
        return
    }

    // Give them back all their login goodies in one package:
    // - Their account information
    // - Their temporary pass (access token)
    // - Their long-term pass (refresh token)
    respondWithJSON(w, http.StatusOK, response{
        User: User{
            ID:        user.ID,        // Their unique ID number
            CreatedAt: user.CreatedAt, // When they first signed up
            UpdatedAt: user.UpdatedAt, // When their info was last changed
            Email:     user.Email,     // Their email address
        },
        Token:        accessToken,  // Their temporary pass
        RefreshToken: refreshToken, // Their long-term pass
    })
}
