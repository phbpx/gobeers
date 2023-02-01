package auth

import (
	"context"
	"fmt"

	"github.com/corabank/go-starter/business/sys/validate"
	"github.com/golang-jwt/jwt/v4"
)

// Claims represents the authorization claims transmitted via a JWT.
type Claims struct {
	jwt.RegisteredClaims

	BusinessID string `json:"business_id" validate:"required,uuid"`
	PersonID   string `json:"person_id" validate:"uuid"`
	AppID      string `json:"azp" validate:"required"`
}

// Validate validates claims.
func (c *Claims) Validate() error {
	if claimsErr := validate.Check(c); claimsErr != nil {
		return fmt.Errorf("invalid user claims: %w", claimsErr)
	}

	return nil
}

// =============================================================================

// ctxKey represents the type of value for the context key.
type ctxKey int

// key is used to store/retrieve a Claims value from a context.Context.
const key ctxKey = 1

// SetClaims stores the claims in the context.
func SetClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, key, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) Claims {
	v, ok := ctx.Value(key).(Claims)
	if !ok {
		return Claims{}
	}
	return v
}
