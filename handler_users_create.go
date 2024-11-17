package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/srinivassivaratri/Chirpy/internal/auth"
	"github.com/srinivassivaratri/Chirpy/internal/database"
)

type User struct {
	ID          uuid.UUID `json:"id"`            // A unique identifier for each user, like a fingerprint. Uses UUID (Universally Unique ID) to avoid duplicates
	CreatedAt   time.Time `json:"created_at"`    // Records when the user first signed up, like a birth certificate date
	UpdatedAt   time.Time `json:"updated_at"`    // Tracks when user info was last changed, like updating your driver's license
	Email       string    `json:"email"`         // User's email address for login and contact, like a digital mailbox
	Password    string    `json:"-"`             // User's password, hidden from JSON output (that's what "-" means) for security
	IsChirpyRed bool      `json:"is_chirpy_red"` // Whether user has premium features (true) or free account (false), like a VIP pass
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:          user.ID,          // Copies the user's unique ID number (like a digital fingerprint) from the database to send back
			CreatedAt:   user.CreatedAt,   // Copies the timestamp of when user first signed up from database to send back
			UpdatedAt:   user.UpdatedAt,   // Copies the timestamp of user's last info update from database to send back
			Email:       user.Email,       // Copies the user's email address from database to send back
			IsChirpyRed: user.IsChirpyRed, // Copies whether user has premium features (true/false) from database to send back
		},
	})
}
