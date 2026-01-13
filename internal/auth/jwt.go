package auth

import "github.com/golang-jwt/jwt/v5"

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
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}

	// NOTE: Use jwt.io to check your test tokens
	return tokenString, nil
}

func (a *JWTAuthenticator) ValidateToken(tokenString string) (*jwt.Token, error) {
	return nil, nil
}
