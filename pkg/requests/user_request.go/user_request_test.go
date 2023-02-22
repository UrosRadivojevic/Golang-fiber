package user_request_test

import (
	"testing"

	"github.com/go-playground/validator"
	"github.com/stretchr/testify/require"
	"github.com/urosradivojevic/health/pkg/requests/user_request.go"
)

func TestShouldReturnErrorWhenEmpty(t *testing.T) {
	// arrange
	t.Parallel()
	assert := require.New(t)
	req := user_request.UserRequest{}
	val := validator.New()

	// act
	err := val.Struct(req).(validator.ValidationErrors)

	// assert
	assert.Error(err)
	assert.Len(err, 2)
}

func TestValide_Success(t *testing.T) {
	// arrange
	t.Parallel()
	assert := require.New(t)
	req := user_request.UserRequest{
		Username:  "uros",
		Password:  "uros1234",
		Firstname: "Uros",
	}
	val := validator.New()
	// act
	err := val.Struct(req)
	// assert
	assert.Nil(err)
}

func TestPartial(t *testing.T) {
	// arrange
	t.Parallel()
	assert := require.New(t)
	req := user_request.UserRequest{
		Username:  "",
		Password:  "uros1234",
		Firstname: "Uros",
	}
	val := validator.New()
	// act
	err := val.Struct(req).(validator.ValidationErrors)
	// assert
	assert.Error(err)
	assert.Len(err, 1)
}
