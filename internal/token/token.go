package token

import (
	"github.com/golang-jwt/jwt/v5"
)

type Token interface {
	GenerateToken(userID string, role string) (string, error)
	Parse(token string) (map[string]interface{}, error)
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

func (j *token) Parse(tokenString string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

