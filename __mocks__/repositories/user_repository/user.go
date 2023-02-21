package user_repository

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/urosradivojevic/health/pkg/model"
	"github.com/urosradivojevic/health/pkg/repositories/user_repository"
)

type MockRepo struct {
	mock.Mock
}

var _ user_repository.Interface = &MockRepo{}

func (u *MockRepo) GetByUsername(ctx context.Context, username string) (model.User, error) {
	args := u.Called(ctx, username)
	return args.Get(0).(model.User), args.Error(1)
}

func (u *MockRepo) Register(ctx context.Context, user model.User) error {
	args := u.Called(ctx, user)
	return args.Error(0)
}

func (u *MockRepo) Exists(ctx context.Context, username string) bool {
	args := u.Called(ctx, username)
	return args.Bool(0)
}
