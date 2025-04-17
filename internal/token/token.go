package token

import (
	"github.com/golang-jwt/jwt/v5"
)

type Token interface {
	GenerateToken(userID string, role string) (string, error)
}

type token struct {
	secret []byte
}

func NewToken(secret string) Token {
	return &token{
		secret: []byte(secret),
	}
}

func (j *token) GenerateToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

