// Package auth provides authentication and authorization support.
package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// jwt parser
var parser = jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}))

// Authenticate processes the token to validate the sender's token is valid.
func Authenticate(ctx context.Context, bearerToken string) (Claims, error) {

	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return Claims{}, errors.New("expected authorization header format: Bearer <token>")
	}

	var claims Claims
	token, _, err := parser.ParseUnverified(parts[1], &claims)
	if err != nil {
		return Claims{}, fmt.Errorf("error parsing token: %w", err)
	}

	// Perform an extra level of token verification.

	// Validate KID

	kidRaw, exists := token.Header["kid"]
	if !exists {
		return Claims{}, fmt.Errorf("kid missing from header: %w", err)
	}

	if _, ok := kidRaw.(string); !ok {
		return Claims{}, fmt.Errorf("kid malformed: %w", err)
	}

	// Validate claims

	if claimsErr := claims.Validate(); claimsErr != nil {
		return Claims{}, fmt.Errorf("invalid claims: %w", claimsErr)
	}

	return claims, nil
}
