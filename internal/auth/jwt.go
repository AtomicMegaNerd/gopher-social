package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthenticator struct {
	secret   string // WARNING: Private key, do not expose
	audience string
	issuer   string
}

func NewJWTAuthenticator(secret, audience, issuer string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secret:   secret,
		audience: audience,
		issuer:   issuer,
	}
}

func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Make sure the secret is configured otherwise we might sign tokens with an empty key
	if a.secret == "" {
		return "", errors.New("jwt secret is not configured")
	}

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	// NOTE: Use jwt.io to check your test tokens
	return tokenString, nil
}

// NOTE: See https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries to avoid
// making critical security errors with JWT!
//
// The none algorithm is a curious addition to JWT. It is intended to be used for situations where
// the integrity of the token has already been verified. Interestingly enough, it is one of only
// two algorithms that are mandatory to implement (the other being HS256).
// Unfortunately, some libraries treated tokens signed with the none algorithm as a valid token
// with a verified signature. The result? Anyone can create their own "signed" tokens with whatever
// payload they want, allowing arbitrary account access on some systems.
//
// Putting together such a token is easy. Modify the above example header to contain "alg": "none"
// instead of HS256. Make any desired changes to the payload. Use an empty signature (i.e.
// signature = ""). Most (hopefully all?) implementations now have a basic check to prevent this
// attack: if a secret key was provided, then token verification will fail for tokens using the
// none algorithm. This is a good idea, but it doesn't solve the underlying problem: attackers
// control the choice of algorithm. Let's keep digging.
func (a *JWTAuthenticator) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.audience),
		jwt.WithIssuer(a.audience),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
