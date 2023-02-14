package handler_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/handler"
	"github.com/urosradivojevic/health/pkg/model"
)

func TestGetMovie_Success(t *testing.T) {
	// Arrange
	// t.Parallel()
	data := struct {
		ID       string
		Movie    string
		Watched  bool
		Year     int
		LeadRole string
	}{
		ID:       "63ea5a4bddc25714e92bfc1e",
		Movie:    "DieHard",
		Watched:  true,
		Year:     2005,
		LeadRole: "Bruce",
	}
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Get("/movie/:id", handler.GetMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	req := httptest.NewRequest(http.MethodGet, "/movie/63ea5a4bddc25714e92bfc1e", nil)

	// Act
	res, _ := app.Test(req)
	// when  kad ti dam ovo rezulat funkcija(ovo)
	// given result
	// then ocekujem da mi je rezultat to

	movie := model.Netflix{}
	bytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	_ = json.Unmarshal(bytes, &movie)

	// Assert
	assert.Equal(data.Movie, movie.Movie)
	assert.Equal(data.ID, movie.ID.Hex())
	assert.Equal(data.LeadRole, movie.LeadRole)
	assert.Equal(data.Year, movie.Year)
	assert.Equal(data.Watched, movie.Watched)
	assert.NoError(err)
	assert.Equal(fiber.StatusOK, res.StatusCode)
}

func TestGetMovies_InvalidObjectID(t *testing.T) {
	//arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Get("/movie/:id", handler.GetMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	req := httptest.NewRequest(http.MethodGet, "/movie/63e684f9afd8b30f56511", nil)

	//act

	res, _ := app.Test(req)

	//Assert

	assert.Equal(fiber.StatusBadRequest, res.StatusCode)

}

func TestCreateMovie_InvalidEntity(t *testing.T) {
	//arrange

	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Post("/movie", handler.CreateMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))

	body := []byte(`{
		aaaaa
	}`)
	req := httptest.NewRequest(http.MethodPost, "/movie", bytes.NewBuffer(body))
	//Act
	res, err := app.Test(req)

	//Assert
	assert.NoError(err)
	assert.Equal(fiber.StatusUnprocessableEntity, res.StatusCode)

}

func TestCreateMovie_Success(t *testing.T) {
	//arrange
	data := struct {
		Movie    string `json:"movie"`
		Watched  bool   `json:"watched"`
		Year     int    `json:"year"`
		LeadRole string `json:"leadrole"`
	}{
		Movie:    "DieHard",
		Watched:  true,
		Year:     2005,
		LeadRole: "Bruce",
	}
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Post("/movie", handler.CreateMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	// body := []byte(`{
	// 	"movie": "DieHard",
	// 	"watched": true,
	// 	"year": 2005,
	// 	"leadrole": "Bruce"
	// }`)
	b, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/movie", bytes.NewBuffer(b))
	//Act
	res, err := app.Test(req)

	movie := model.Netflix{}
	bytes, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	_ = json.Unmarshal(bytes, &movie)

	//Assert

	assert.NoError(err)
	assert.Equal(fiber.StatusCreated, res.StatusCode)
	// assert.Equal("DieHard", movie.Movie)
	// assert.Equal(2005, movie.Year)
	// assert.Equal(true, movie.Watched)
	// assert.Equal("Bruce", movie.LeadRole)
}

func TestDeleteMovie_Success(t *testing.T) {
	//arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Delete("/movie/:id", handler.DeleteMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	req := httptest.NewRequest(http.MethodDelete, "/movie/63e684f9afd8b30f56511e46", nil)

	//act
	res, err := app.Test(req)

	//assert
	assert.NoError(err)
	assert.Equal(fiber.StatusNoContent, res.StatusCode)
}

func TestDeleteMovie_InvalidObjectID(t *testing.T) {
	//arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Delete("/movie/:id", handler.DeleteMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	req := httptest.NewRequest(http.MethodDelete, "/movie/63e684f9afd8b30f", nil)

	//act

	res, _ := app.Test(req)

	//Assert

	assert.Equal(fiber.StatusBadRequest, res.StatusCode)
}

func TestGetMovies_Success(t *testing.T) {
	//arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Get("/movies", handler.GetMovies(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)

	//act
	res, err := app.Test(req)

	//assert
	assert.NoError(err)
	assert.Equal(fiber.StatusOK, res.StatusCode)
}

func TestMarkedAsWatched_InvalidObjectID(t *testing.T) {
	//arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Put("/movie/:id", handler.MarkAsWatched(c.GetNetflixRepository()))
	req := httptest.NewRequest(http.MethodPut, "/movie/63e684f9afd8b30f", nil)

	//act

	res, _ := app.Test(req)

	//Assert

	assert.Equal(fiber.StatusBadRequest, res.StatusCode)
}

func TestMarkedAsWatched_Success(t *testing.T) {
	//arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Put("/movie/:id", handler.MarkAsWatched(c.GetNetflixRepository()))
	req := httptest.NewRequest(http.MethodPut, "/movie/63e22d6a22d15a2f27e6ba70", nil)

	//act
	res, err := app.Test(req)

	//assert
	assert.NoError(err)
	assert.Equal(fiber.StatusNoContent, res.StatusCode)

}
