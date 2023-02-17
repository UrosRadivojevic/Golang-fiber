package cache_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urosradivojevic/health/pkg/cache"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSetMovie_Success(t *testing.T) {
	// arragne
	assert := require.New(t)
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	cacheReposisotry := cache.New(c.GetRedisClient())
	movie := model.Netflix{
		ID:       primitive.NewObjectID(),
		Movie:    "How High",
		Watched:  true,
		Year:     2001,
		LeadRole: "Method Man and Redman",
	}
	// act
	err := cacheReposisotry.SetMovie(context.Background(), movie)
	movie1, _ := cacheReposisotry.Get(context.Background(), movie.ID.Hex())
	// assert
	assert.NoError(err)
	assert.Equal(movie1.ID.Hex(), movie.ID.Hex())
	assert.Equal(movie1.Watched, movie.Watched)
	assert.Equal(movie1.Year, movie.Year)
	assert.Equal(movie1.LeadRole, movie.LeadRole)
	assert.Equal(movie1.Movie, movie.Movie)
}

func TestGet_Success(t *testing.T) {
	// arragne
	assert := require.New(t)
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	cacheReposisotry := cache.New(c.GetRedisClient())
	movie := model.Netflix{
		ID:       primitive.NewObjectID(),
		Movie:    "How High",
		Watched:  true,
		Year:     2001,
		LeadRole: "Method Man and Redman",
	}
	err := cacheReposisotry.SetMovie(context.Background(), movie)

	// act
	movie1, _ := cacheReposisotry.Get(context.Background(), movie.ID.Hex())

	// assert
	assert.NoError(err)
	assert.Equal(movie1.ID.Hex(), movie.ID.Hex())
	assert.Equal(movie1.Watched, movie.Watched)
	assert.Equal(movie1.Year, movie.Year)
	assert.Equal(movie1.LeadRole, movie.LeadRole)
	assert.Equal(movie1.Movie, movie.Movie)
}

func TestDelete_Success(t *testing.T) {
	// arragne
	assert := require.New(t)
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	cacheReposisotry := cache.New(c.GetRedisClient())
	movie := model.Netflix{
		ID:       primitive.NewObjectID(),
		Movie:    "How High",
		Watched:  true,
		Year:     2001,
		LeadRole: "Method Man and Redman",
	}
	_ = cacheReposisotry.SetMovie(context.Background(), movie)

	// act
	err := cacheReposisotry.Delete(context.Background(), movie.ID.Hex())
	m, err1 := cacheReposisotry.Get(context.Background(), movie.ID.Hex())

	// assert
	assert.NoError(err)
	assert.Error(err1)
	assert.Empty(m)
}
