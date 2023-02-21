package login

import (
	"context"
	"errors"

	"github.com/urosradivojevic/health/pkg/model"
	"github.com/urosradivojevic/health/pkg/repositories/user_repository"
	"github.com/urosradivojevic/health/pkg/services/hasher"
)

var (
	ErrUsernameNotFound   = errors.New("Username not found.")
	ErrInvalidCredentials = errors.New("Username & Password combination incorect!")
)

type Interface interface {
	Login(username string, password string) (model.User, error)
}

type LoginService struct {
	repo user_repository.Interface
	hash hasher.Interface
}

func New(repo user_repository.Interface, hash hasher.Interface) LoginService {
	return LoginService{
		repo: repo,
		hash: hash,
	}
}

func (l LoginService) Login(username string, password string) (model.User, error) {
	user, err := l.repo.GetByUsername(context.Background(), username)
	if err != nil {
		return model.User{}, ErrUsernameNotFound
	}
	if !l.hash.Verify(password, []byte(user.Password)) {
		return model.User{}, ErrInvalidCredentials
	}
	return user, nil
}
