package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/srinivassivaratri/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	// Verify API key first
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing API key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key", errors.New("invalid api key"))
		return
	}

	type webhookBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	payload := webhookBody{}
	err = decoder.Decode(&payload)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// If it's not a user.upgraded event, return 204 immediately
	if payload.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Parse the user ID from string to UUID
	userID, err := uuid.Parse(payload.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID format", err)
		return
	}

	// Attempt to upgrade the user to Chirpy Red
	_, err = cfg.db.UpgradeToChirpyRed(r.Context(), userID)
	if err != nil {
		// If user not found, return 404
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	// Success - return 204
	w.WriteHeader(http.StatusNoContent)
}
