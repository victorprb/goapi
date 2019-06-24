package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	app.writeJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

// The clientError helper writes a client error response to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	app.writeJSON(w, status, http.StatusText(status))
}

// The notFound helper writes a not found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// The writeJSON helper writes a JSON response to the user
func (app *application) writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// authenticateJWTToken verify the JWT token from a request and it returns the claims
func authenticateJWTToken(secretKey string, r *http.Request) (map[string]interface{}, error) {
	jwtToken, err := extractJWTToken(r)

	if err != nil {
		return nil, fmt.Errorf("Failed get JWT token")
	}

	claims, err := parseJWT(jwtToken, secretKey)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse token")
	}

	return claims, nil
}

// extractJWTToken extracts bearer token from Authorization header
func extractJWTToken(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")

	if tokenString == "" {
		return "", fmt.Errorf("Could not find token")
	}

	tokenString, err := stripTokenPrefix(tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Strips 'Token' or 'Bearer' prefix from token string
func stripTokenPrefix(tok string) (string, error) {
	// split token to 2 parts
	tokenParts := strings.Split(tok, " ")

	if len(tokenParts) < 2 {
		return tokenParts[0], nil
	}

	return tokenParts[1], nil
}

// parseJWT parses a JWT and returns Claims object
func parseJWT(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate if the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// secretKey is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secretKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}