package main

import (
	"net/http"
	"time"

	"github.com/srinivassivaratri/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	// Define what our response JSON will look like - just a token field that will be labeled "token" in JSON
	type response struct {
		Token string `json:"token"`
	}

	// Look in the request headers for "Authorization: Bearer <token>" and grab just the token part
	// This is like checking someone's ID at the door of a club
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing refresh token", err)
		return
	}

	// Take that refresh token and look up which user it belongs to in our database
	// Like using a coat check ticket to get back your coat
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	// Make a fresh JWT access token for this user that expires in 1 hour
	// Like giving someone a new temporary visitor badge after they show proper ID
	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}

	// Package up the new access token in our response format and send it back
	// Like handing over the new visitor badge to the person
	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
} 