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

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "pursuit-access"
)

// ErrNoAuthHeaderIncluded
var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")

// HasPassword -
func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// CheckPasswordHash -
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// MakeJWT token
func MakeJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {

	signingKey := []byte(tokenSecret) // Convert the token secret to a byte slice

	if len(signingKey) == 0 { // Check if the signing key is empty
		return "", errors.New("empty signing key")
	}
	if userID == uuid.Nil { // Check if the user ID is nil
		return "", errors.New("empty user ID")
	}
	if expiresIn == 0 { // Check if the expiration duration is zero
		return "", errors.New("empty expiration duration")
	}
	// Create a new JWT Token with specified signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString(signingKey) // Sign the token with the signing key and return it
}

// Validate JWT -
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claimsStruct := jwt.RegisteredClaims{} // Create a new instance of the RegisteredClaims struct
	token, err := jwt.ParseWithClaims(     // Parse the token string and validate it using the provided secret
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return uuid.Nil, err
	}

	userIDString, err := token.Claims.GetSubject() // Get the subject (user ID) from the token claims
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer() // Get the issuer from the token claims
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) { // Check if the issuer matches the expected token type
		return uuid.Nil, errors.New("invalid token issuer")
	}

	id, err := uuid.Parse(userIDString) // Parse the user ID string into a UUID
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID: %w", err) // Return an error if the user ID is invalid
	}
	return id, nil // Return the parsed user ID
}

// GetBearerToken
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization") // Get the Authorization header from the request headers
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ") // Split the header into parts
	// Check if the header is malformed
	// The first part should be "Bearer" and the second part should be the token
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil // Return the token part of the header
}

// MakeRefreshToken makes a random 256 bit token encoded in hex
func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)  // Create a byte slice of length 32
	_, err := rand.Read(token) // Fill the byte slice with random data
	if err != nil {
		return "", err // Return an error if the random data generation fails
	}
	return hex.EncodeToString(token), nil // Encode the byte slice to a hex string and return it
}
