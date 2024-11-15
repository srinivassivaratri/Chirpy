package main

import (
	"net/http"

	"github.com/srinivassivaratri/Chirpy/internal/auth"
)

// This function handles revoking (invalidating) refresh tokens. It's attached to the apiConfig struct
// so it has access to the database and other configuration
func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// This line extracts the refresh token from the Authorization header in the HTTP request
	// The token is expected to be in the format "Bearer <token>"
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		// If there's no token or it's in the wrong format, tell the user they need to provide a valid token
		// 401 Unauthorized is the standard response for missing/invalid auth
		respondWithError(w, http.StatusUnauthorized, "Missing refresh token", err)
		return
	}

	// This calls the database to mark the token as revoked by setting its revoked_at timestamp
	// We pass the request context in case we need to cancel the operation
	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		// If the database operation fails for any reason, we return a 500 server error
		// This could happen if the database is down or there's a connection issue
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}

	// If everything worked, we send back 204 No Content
	// 204 is used when the operation succeeded but there's no data to return
	// It's common for DELETE-like operations
	w.WriteHeader(http.StatusNoContent)
} 