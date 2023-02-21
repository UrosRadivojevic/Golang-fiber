package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

type (
	Interface interface {
		Hash(password string) ([]byte, error)
		Verify(password string, hashPassword []byte) bool
	}
	Hasher struct {
		cost int32
	}
)

func New(cost int32) *Hasher {
	return &Hasher{
		cost: cost,
	}
}

func (h *Hasher) Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), int(h.cost))
}

func (h *Hasher) Verify(password string, hashPassword []byte) bool {
	return bcrypt.CompareHashAndPassword(hashPassword, []byte(password)) == nil
}
