package main

import (
	"net/http"
	"time"

	"github.com/srinivassivaratri/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
    // This function runs whenever someone wants to get a new access token using their refresh token
    // cfg contains our app's settings and tools
    // w lets us send back a response to the user
    // r contains all the details about their request to us

    type response struct {
        Token string `json:"token"`
    }
    // We create a template for what we'll send back
    // It will be a piece of data with just one field called "token"
    // The json:"token" part tells Go to name this field "token" when converting to JSON format

    refreshToken, err := auth.GetBearerToken(r.Header)
    // We look in the request headers for something like "Bearer abc123..."
    // We extract just the token part (abc123...)
    // If we find it, it goes in refreshToken. If something goes wrong, err will tell us what

    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Missing refresh token", err)
        return
    }
    // If we couldn't find a valid token in the headers:
    // - Tell the user they're not authorized (401 error)
    // - Include a message explaining why
    // - Stop processing this request

    user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
    // We check our database to see if this refresh token is valid
    // If it is, we get back the user account it belongs to
    // If not, err will tell us what went wrong

    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
        return
    }
    // If we couldn't find a valid user for this token:
    // - Tell them they're not authorized (401 error)
    // - Explain their token isn't valid
    // - Stop processing this request

    accessToken, err := auth.MakeJWT(
        user.ID,        // Include the user's ID in the token
        cfg.jwtSecret,  // Use our secret key to sign it
        time.Hour,      // Make it expire in 1 hour
    )
    // Create a new temporary access token
    // It's like a hall pass that proves who they are
    // We sign it with our secret so we know we created it
    // It will stop working after an hour for security

    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
        return
    }
    // If something went wrong making the new token:
    // - Tell them there was a server error (500)
    // - Explain we couldn't create their token
    // - Stop processing this request

    respondWithJSON(w, http.StatusOK, response{
        Token: accessToken,
    })
    // Everything worked! So we:
    // - Package up the new access token in our response template
    // - Convert it to JSON format
    // - Send it back with a success status (200 OK)
}