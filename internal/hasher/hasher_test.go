//go:build unit

package hasher


import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasher_HashAndEqual_Success(t *testing.T) {
	h := NewHasher()
	password := "myStrongPassword123"

	hashed, err := h.Hash(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)

	match := h.Equal(hashed, password)
	assert.True(t, match)
}

func TestHasher_Equal_InvalidPassword(t *testing.T) {
	h := NewHasher()

	password := "correctPassword"
	wrongPassword := "wrongPassword"

	hashed, err := h.Hash(password)
	assert.NoError(t, err)

	match := h.Equal(hashed, wrongPassword)
	assert.False(t, match)
}

func TestHasher_Equal_InvalidHashFormat(t *testing.T) {
	h := NewHasher()

	invalidHash := "this_is_not_a_valid_bcrypt_hash"
	password := "anyPassword"

	match := h.Equal(invalidHash, password)
	assert.False(t, match)
}

func TestHasher_Hash_ErrorHandling(t *testing.T) {
	h := NewHasher()

	hashed, err := h.Hash("")
	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)

	match := h.Equal(hashed, "")
	assert.True(t, match)
}


