package requests_test

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"github.com/urosradivojevic/health/pkg/requests"
)

func TestShouldReturnErrorsWhenEmpty(t *testing.T) {
	// arrange
	t.Parallel()
	assert := require.New(t)
	req := requests.CreateMovieRequest{}
	val := validator.New()

	// act
	err := val.Struct(req).(validator.ValidationErrors)

	// assert
	assert.Error(err)
	assert.Len(err, 4)
}

func TestValide_Success(t *testing.T) {
	// arrange
	t.Parallel()
	assert := require.New(t)
	req := requests.CreateMovieRequest{
		Movie:    "Die Hard",
		Watched:  true,
		Year:     2005,
		LeadRole: "Bruce",
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
	req := requests.CreateMovieRequest{
		Movie:    "",
		Watched:  true,
		Year:     2005,
		LeadRole: "",
	}
	val := validator.New()
	// act
	err := val.Struct(req).(validator.ValidationErrors)
	// assert
	assert.Error(err)
	assert.Len(err, 2)
}
