package movie_repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/repositories/movie_repository"
	"github.com/urosradivojevic/health/pkg/requests"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInsertOneMovie_Success(t *testing.T) {
	// arrange
	// t.Parallel()
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	assert := require.New(t)
	c := container.New("testing")
	col := c.GetMongoCollection("watchlist")
	movieRepository := movie_repository.New(col)
	movie := requests.CreateMovieRequest{
		Movie:    "How High",
		Watched:  true,
		Year:     2001,
		LeadRole: "Method Man and Redman",
	}

	// act
	id, err := movieRepository.InsertOneMovie(movie)
	movie1, err1 := movieRepository.GetOneMovie(id.Hex())

	// assert
	assert.NoError(err1)
	assert.NotEmpty(id)
	assert.NoError(err)
	assert.Len(id.Hex(), 24)
	assert.Equal(movie1.ID.Hex(), id.Hex())
	assert.Equal(movie1.Movie, movie.Movie)
	assert.Equal(movie1.LeadRole, movie.LeadRole)
	assert.Equal(movie1.Watched, movie.Watched)
	assert.Equal(movie1.Year, movie.Year)
}

func TestUpdateOneMovie_Success(t *testing.T) {
	// arrange
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	assert := require.New(t)
	c := container.New("testing")
	col := c.GetMongoCollection("watchlist")
	movieRepository := movie_repository.New(col)
	movie := requests.CreateMovieRequest{
		Movie:    "How High",
		Watched:  false,
		Year:     2001,
		LeadRole: "Method Man and Redman",
	}
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})
	id, _ := movieRepository.InsertOneMovie(movie)

	// act
	err := movieRepository.UpdateOneMovie(id.Hex())
	movie1, _ := movieRepository.GetOneMovie(id.Hex())

	// assert
	assert.NoError(err)
	assert.True(movie1.Watched)
}

func TestDeleteOneMovie_Success(t *testing.T) {
	// arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	c := container.New("testing")
	col := c.GetMongoCollection("watchlist")
	movieRepository := movie_repository.New(col)
	movie := requests.CreateMovieRequest{
		Movie:    "How High",
		Watched:  false,
		Year:     2001,
		LeadRole: "Method Man and Redman",
	}
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})
	id, _ := movieRepository.InsertOneMovie(movie)

	// act
	err := movieRepository.DeleteOneMovie(id.Hex())
	movies, _ := movieRepository.GetAllMovies()

	// assert
	assert.NoError(err)
	assert.Empty(movies)
}

func TestGetAllMovies_Success(t *testing.T) {
	// arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	c := container.New("testing")
	col := c.GetMongoCollection("watchlist")
	movieRepository := movie_repository.New(col)
	movie := requests.CreateMovieRequest{
		Movie:    "How High",
		Watched:  false,
		Year:     2001,
		LeadRole: "Method Man and Redman",
	}
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})
	id, _ := movieRepository.InsertOneMovie(movie)

	// act
	movies, err := movieRepository.GetAllMovies()

	// assert
	assert.NoError(err)
	assert.Equal(movies[0].ID.Hex(), id.Hex())
	assert.Equal(movies[0].Movie, movie.Movie)
	assert.Equal(movies[0].LeadRole, movie.LeadRole)
	assert.Equal(movies[0].Watched, movie.Watched)
	assert.Equal(movies[0].Year, movie.Year)
}
