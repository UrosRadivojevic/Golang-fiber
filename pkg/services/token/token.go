package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Interface interface {
	Generate() (string, error)
}

type Service struct{}

func New() Service {
	return Service{}
}

func (s Service) Generate() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
