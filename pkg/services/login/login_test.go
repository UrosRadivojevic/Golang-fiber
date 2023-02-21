package login_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/urosradivojevic/health/__mocks__/repositories/user"
	"github.com/urosradivojevic/health/__mocks__/services/hasher"
	"github.com/urosradivojevic/health/pkg/model"
	"github.com/urosradivojevic/health/pkg/services/login"
)

func TestLoginUsernameNotFound(t *testing.T) {
	// arrange
	assert := require.New(t)
	userRepo := &user.MockRepo{}
	hasherService := &hasher.HasherMock{}
	l := login.New(userRepo, hasherService)
	userRepo.On("GetByUsername", context.Background(), "testUsername").Return(model.User{Firstname: "Mare"}, errors.New("Test error"))
	// act
	user, err := l.Login("testUsername", "test")

	// assert
	assert.ErrorIs(err, login.ErrUsernameNotFound)
	assert.Empty(user)
	userRepo.AssertExpectations(t)
	hasherService.AssertExpectations(t)
	hasherService.AssertNotCalled(t, "Verify")
}

func TestLoginInvalidCredentials(t *testing.T) {
	// arrange
	assert := require.New(t)
	userRepo := &user.MockRepo{}
	hasherService := &hasher.HasherMock{}
	l := login.New(userRepo, hasherService)
	userRepo.On("GetByUsername", context.Background(), "testUsername").Return(model.User{Firstname: "KD"}, nil)
	hasherService.On("Verify", mock.Anything, mock.Anything).Return(false)
	// act
	user, err := l.Login("testUsername", "test")
	// assert
	assert.ErrorIs(err, login.ErrInvalidCredentials)
	assert.Empty(user)
	userRepo.AssertExpectations(t)
	userRepo.AssertCalled(t, "GetByUsername", context.Background(), "testUsername")
	hasherService.AssertExpectations(t)
	hasherService.AssertCalled(t, "Verify", mock.Anything, mock.Anything)
}

func TestLoginSuccess(t *testing.T) {
	// arrange
	assert := require.New(t)
	userRepo := &user.MockRepo{}
	hasherService := &hasher.HasherMock{}
	l := login.New(userRepo, hasherService)
	userRepo.On("GetByUsername", context.Background(), "testUsername").Return(model.User{Firstname: "KD"}, nil)
	hasherService.On("Verify", mock.Anything, mock.Anything).Return(true)
	// act
	user, err := l.Login("testUsername", "test")
	// assert
	assert.NoError(err)
	assert.NotEmpty(user)
	userRepo.AssertExpectations(t)
	userRepo.AssertCalled(t, "GetByUsername", context.Background(), "testUsername")
	hasherService.AssertExpectations(t)
	hasherService.AssertCalled(t, "Verify", mock.Anything, mock.Anything)
}
