//go:build unit

package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestToken_GenerateAndParse_Success(t *testing.T) {
	secret := "mysecret"
	tk := NewToken(secret)

	userID := "user123"
	role := "Client"

	tokenStr, err := tk.GenerateToken(userID, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	claims, err := tk.Parse(tokenStr)
	assert.NoError(t, err)

	assert.Equal(t, userID, claims["user_id"])
	assert.Equal(t, role, claims["role"])
}

func TestToken_Parse_InvalidSignature(t *testing.T) {
	// Создаём токен с другим секретом
	secret := "correct-secret"
	wrongSecret := "wrong-secret"
	tk := NewToken(secret)

	// Токен подписан другим ключом
	foreign := NewToken(wrongSecret)
	tokenStr, err := foreign.GenerateToken("user123", "Client")
	assert.NoError(t, err)

	claims, err := tk.Parse(tokenStr)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestToken_Parse_InvalidMethod(t *testing.T) {
	secret := "mysecret"

	claims := jwt.MapClaims{
		"user_id": "user123",
		"role":    "Client",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenStr, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	assert.NoError(t, err)

	tk := NewToken(secret)
	parsedClaims, err := tk.Parse(tokenStr)
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
}

func TestToken_Parse_MalformedToken(t *testing.T) {
	secret := "mysecret"
	tk := NewToken(secret)

	// Нарушаем формат токена
	malformed := "this.is.not.a.jwt"

	claims, err := tk.Parse(malformed)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestToken_Parse_EmptyToken(t *testing.T) {
	tk := NewToken("secret")

	claims, err := tk.Parse("")
	assert.Error(t, err)
	assert.Nil(t, claims)
}