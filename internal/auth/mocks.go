package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TestAuthenticator struct {
}

const secret = "test"

var testClaims = jwt.MapClaims{
	"aud": "test-audience",
	"iss": "test-audience",
	"sub": int64(1),
	"exp": time.Now().Add(time.Hour).Unix(),
}

// WARNING: Do not use this method in production it is for unit tests only
func (t *TestAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	// Ignore the actual claims and use the fake ones above
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	return token.SignedString([]byte(secret))
}

// WARNING: Do not use this method in production it is for unit tests only
func (t *TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
}
