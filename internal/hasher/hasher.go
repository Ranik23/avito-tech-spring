package hasher

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	Hash(password string) (string, error)
	Equal(hashedPassword string, password string) bool
}

type hasher struct {}

func NewHasher() Hasher {
	return &hasher{}
}

func (h *hasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (h *hasher) Equal(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
