package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/srinivassivaratri/Chirpy/internal/auth"
	"github.com/srinivassivaratri/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	// This function handles requests to delete a chirp (like deleting a tweet)
	// cfg contains settings and database connection
	// w is where we write our response back to the user
	// r contains all the details about what the user is requesting

	chirpIDString := r.PathValue("chirpID")
	// When someone visits /chirps/123, this grabs the "123" part
	// We store it as chirpIDString since it's still in text form

	chirpID, err := uuid.Parse(chirpIDString)
	// Takes the ID text and converts it to a special format called UUID
	// UUIDs are like super-long random numbers that are guaranteed to be unique
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	// If the ID text isn't in the right format, tell the user they made a mistake

	token, err := auth.GetBearerToken(r.Header)
	// Looks in the request for a secret code (token) that proves who the user is
	// It's like checking someone's ID card at a club
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing access token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	// Makes sure the ID card (token) is real and not fake
	// If it's real, we get back the user's unique identifier
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	// Asks our database: "Do you have a chirp with this ID? Show me its details"
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can only delete your own chirps", errors.New("unauthorized deletion attempt"))
		return
	}
	// Checks if the person trying to delete the chirp actually created it
	// Like making sure you can only delete your own posts, not someone else's

	err = cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpID,
		UserID: userID,
	})
	// Tells the database: "Please delete this specific chirp"
	// We include who's asking to delete it as a double-check
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	// If everything worked, we send back a simple "done" message
	// StatusNoContent (204) means "success, but I have nothing else to tell you"
}
