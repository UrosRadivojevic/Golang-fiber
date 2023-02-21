package token

import (
	"crypto/rand"
	"fmt"
)

type Interface interface {
	Generate(size int) string
}

type Service struct{}

func New() Service {
	return Service{}
}

func (s Service) Generate(size int) string {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}
