package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/srinivassivaratri/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	authorIDStr := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	// Validate sort order
	if sortOrder != "" && sortOrder != "asc" && sortOrder != "desc" {
		respondWithError(w, http.StatusBadRequest, "Invalid sort order. Use 'asc' or 'desc'", nil)
		return
	}

	// Default to "asc" if not specified
	if sortOrder == "" {
		sortOrder = "asc"
	}

	var dbChirps []database.Chirp
	var err error

	// Get chirps with appropriate sort order
	if sortOrder == "desc" {
		dbChirps, err = cfg.db.GetChirpsDesc(r.Context())
	} else {
		dbChirps, err = cfg.db.GetChirpsAsc(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	// Filter by author if specified
	chirps := []Chirp{}
	if authorIDStr != "" {
		authorID, err := uuid.Parse(authorIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID format", err)
			return
		}

		// Only include chirps from specified author
		for _, dbChirp := range dbChirps {
			if dbChirp.UserID == authorID {
				chirps = append(chirps, Chirp{
					ID:        dbChirp.ID,
					CreatedAt: dbChirp.CreatedAt,
					UpdatedAt: dbChirp.UpdatedAt,
					UserID:    dbChirp.UserID,
					Body:      dbChirp.Body,
				})
			}
		}
	} else {
		// Include all chirps
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:        dbChirp.ID,
				CreatedAt: dbChirp.CreatedAt,
				UpdatedAt: dbChirp.UpdatedAt,
				UserID:    dbChirp.UserID,
				Body:      dbChirp.Body,
			})
		}
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
