package token

import (
	"log/slog"
	"github.com/golang-jwt/jwt/v5"
)

type Token interface {
	GenerateToken(userID string, role string) (string, error)
	Parse(token string) (map[string]interface{}, error)
}

type token struct {
	logger *slog.Logger
	secret []byte
}

func NewToken(secret string, logger *slog.Logger) Token {
	return &token{
		secret: []byte(secret),
		logger: logger,
	}
}

func (j *token) GenerateToken(userID string, role string) (string, error) {
	j.logger.Info("Generating token", slog.String("user_id", userID), slog.String("role", role))

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)

	if err != nil {
		j.logger.Error("Error generating token", slog.String("error", err.Error()))
		return "", err
	}

	j.logger.Info("Token generated successfully", slog.String("user_id", userID))
	return tokenString, nil
}

func (j *token) Parse(tokenString string) (map[string]interface{}, error) {
	j.logger.Info("Parsing token", "token", tokenString)

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			j.logger.Error("Invalid signature method", slog.String("method", token.Method.Alg()))
			return nil, jwt.ErrSignatureInvalid
		}
		return j.secret, nil
	})

	if err != nil {
		j.logger.Error("Error parsing token", "error", err)
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		j.logger.Info("Token parsed successfully", slog.String("user_id", claims["user_id"].(string)))
		return claims, nil
	}

	j.logger.Error("Invalid token signature")
	return nil, jwt.ErrSignatureInvalid
}
