package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TokenType is a custom type that represents what kind of authentication token we're dealing with
// It's just a string under the hood, but creating a custom type helps prevent mixing it up with regular strings
type TokenType string

const (
	// TokenTypeAccess defines the name of our access token type
	// We use this constant value "chirpy-access" to identify tokens created by our system
	TokenTypeAccess TokenType = "chirpy-access"
)

// ErrNoAuthHeaderIncluded is an error we return when someone tries to access a protected route
// without including their authentication token in the request header
var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

// ErrNoAPIKeyIncluded is an error we return when someone tries to access a protected route
// without including their API key in the request header
var ErrNoAPIKeyIncluded = errors.New("no api key included in request")

// HashPassword takes a plain text password and turns it into a scrambled version
// This is important because we never want to store actual passwords - only their scrambled form
// That way if someone hacks our database, they can't steal passwords
func HashPassword(password string) (string, error) {
	// Convert password to bytes (computers work with bytes, not strings)
	// Use bcrypt to scramble it with a default level of complexity
	// More complexity = harder to crack but slower to generate
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// If something goes wrong, return empty string and the error
		return "", err
	}
	// Convert the scrambled bytes back to a string so we can store it
	return string(dat), nil
}

// CheckPasswordHash compares a plain text password with a previously hashed (scrambled) version
// This is how we check if someone typed their password correctly during login
// We can't unscramble the hash - we can only create a new hash and compare them
func CheckPasswordHash(password, hash string) error {
	// bcrypt does the comparison for us and returns an error if they don't match
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// MakeJWT creates a signed JWT token for authentication
// Function takes 3 inputs:
// - userID: unique identifier for the user
// - tokenSecret: private key used to sign the token
// - expiresIn: how long the token should be valid for
func MakeJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	// Convert the secret string into bytes so it can be used for signing
	signingKey := []byte(tokenSecret)

	// Create a new JWT token with:
	// - HS256 algorithm for signing
	// - Standard JWT claims containing:
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		// Who created the token (our auth system)
		Issuer: string(TokenTypeAccess),
		// When the token was created (current UTC time)
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		// When the token expires (current time + duration)
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		// Who the token is for (the user's ID)
		Subject: userID.String(),
	})
	// Sign the token with our secret key and return the final string
	return token.SignedString(signingKey)
}

// ValidateJWT takes a JWT token string and secret, validates it, and returns the user ID
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Create an empty struct to store the standard JWT claims (like subject, issuer, expiry)
	claimsStruct := jwt.RegisteredClaims{}

	// Parse the JWT string into a token object, using our secret to verify the signature
	// The anonymous function provides the key for verification
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.Nil, err // Return empty UUID if parsing fails
	}

	// Extract the subject claim (user ID) from the token
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	// Extract the issuer claim to verify this token was issued by our system
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	// Check if the issuer matches our expected token type
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	// Convert the user ID string into a UUID object
	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err)
	}
	return id, nil // Return the valid user ID
}

// GetBearerToken extracts a JWT token from HTTP headers
func GetBearerToken(headers http.Header) (string, error) {
	// Get the Authorization header value from the HTTP request headers
	// This header should contain something like "Bearer eyJ0eXAiOiJKV1QiLC..."
	authHeader := headers.Get("Authorization")

	// If no Authorization header was found, return an error
	// We can't authenticate the request without this header
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	// Split the header value into parts using space as separator
	// A valid header looks like: "Bearer <token>"
	// So splitting "Bearer abc123" gives us ["Bearer", "abc123"]
	splitAuth := strings.Split(authHeader, " ")

	// Check if:
	// 1. We have at least 2 parts (Bearer + token)
	// 2. First part is exactly "Bearer"
	// If either check fails, the header format is wrong
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	// Return just the token part (second element in the split array)
	// For "Bearer abc123", returns "abc123"
	return splitAuth[1], nil
}

// MakeRefreshToken creates a random token used to get new access tokens
func MakeRefreshToken() (string, error) {
	// We need a place to store random data - 32 bytes gives us 256 bits of randomness
	// which makes it practically impossible to guess the token
	randomBytes := make([]byte, 32)

	// We use the computer's built-in random number generator (usually based on
	// electrical noise or timing variations) to fill our byte array with unpredictable values
	_, err := rand.Read(randomBytes)
	if err != nil {
		// If the random generator fails (very rare), we need to let the caller know
		return "", fmt.Errorf("could not generate random bytes: %w", err)
	}

	// The random bytes could contain any values (0-255), which might cause problems in text.
	// So we convert them to hexadecimal (0-9,a-f), which is safe to use anywhere
	return hex.EncodeToString(randomBytes), nil
}

// GetAPIKey extracts an API key from HTTP headers
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAPIKeyIncluded
	}

	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) != 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed api key header")
	}

	return splitAuth[1], nil
}
