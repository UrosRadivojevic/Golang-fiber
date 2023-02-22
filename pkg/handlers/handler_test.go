package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/handlers/create_movie_handler"
	"github.com/urosradivojevic/health/pkg/handlers/delete_movie_handler"
	"github.com/urosradivojevic/health/pkg/handlers/get_movie_handler"
	"github.com/urosradivojevic/health/pkg/handlers/get_movies_handler"
	"github.com/urosradivojevic/health/pkg/handlers/login_handler.go"
	"github.com/urosradivojevic/health/pkg/handlers/mark_as_watched_handler"
	"github.com/urosradivojevic/health/pkg/handlers/register_handler"
	"github.com/urosradivojevic/health/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
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
	app.Get("/movie/:id", get_movie_handler.GetMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	req := httptest.NewRequest(http.MethodGet, "/movie/63ea5a4bddc25714e92bfc1e", nil)

	// Act
	res, _ := app.Test(req)
	// when  kad ti dam ovo rezulat funkcija(ovo)
	// given result
	// then ocekujem da mi je rezultat to

	movie := model.Netflix{}
	bytes, err := io.ReadAll(res.Body)
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
	// arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Get("/movie/:id", get_movie_handler.GetMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	req := httptest.NewRequest(http.MethodGet, "/movie/63e684f9afd8b30f56511", nil)

	// act

	res, _ := app.Test(req)

	// Assert

	assert.Equal(fiber.StatusBadRequest, res.StatusCode)
}

func TestCreateMovie_InvalidEntity(t *testing.T) {
	// arrange

	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Post("/movie", create_movie_handler.CreateMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))

	body := []byte(`{
		aaaaa
	}`)
	req := httptest.NewRequest(http.MethodPost, "/movie", bytes.NewBuffer(body))
	// Act
	res, err := app.Test(req)

	// Assert
	assert.NoError(err)
	assert.Equal(fiber.StatusUnprocessableEntity, res.StatusCode)
}

func TestCreateMovie_Success(t *testing.T) {
	// arrange
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
	app.Post("/movie", create_movie_handler.CreateMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	b, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/movie", bytes.NewBuffer(b))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	// Act
	res, err := app.Test(req)

	movie := model.Netflix{}
	bytes, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	_ = json.Unmarshal(bytes, &movie)

	// Assert

	assert.NoError(err)
	assert.Equal(fiber.StatusCreated, res.StatusCode)
	// assert.Equal("DieHard", movie.Movie)
	// assert.Equal(2005, movie.Year)
	// assert.Equal(true, movie.Watched)
	// assert.Equal("Bruce", movie.LeadRole)
}

func TestDeleteMovie_Success(t *testing.T) {
	// arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Delete("/movie/:id", delete_movie_handler.DeleteMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	req := httptest.NewRequest(http.MethodDelete, "/movie/63e684f9afd8b30f56511e46", nil)

	// act
	res, err := app.Test(req)

	// assert
	assert.NoError(err)
	assert.Equal(fiber.StatusNoContent, res.StatusCode)
}

func TestDeleteMovie_InvalidObjectID(t *testing.T) {
	// arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Delete("/movie/:id", delete_movie_handler.DeleteMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	req := httptest.NewRequest(http.MethodDelete, "/movie/63e684f9afd8b30f", nil)

	// act

	res, _ := app.Test(req)

	// Assert

	assert.Equal(fiber.StatusBadRequest, res.StatusCode)
}

func TestGetMovies_Success(t *testing.T) {
	// arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Get("/movies", get_movies_handler.GetMovies(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)

	// act
	res, err := app.Test(req)

	// assert
	assert.NoError(err)
	assert.Equal(fiber.StatusOK, res.StatusCode)
}

func TestMarkedAsWatched_InvalidObjectID(t *testing.T) {
	// arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Put("/movie/:id", mark_as_watched_handler.MarkAsWatched(c.GetNetflixRepository()))
	req := httptest.NewRequest(http.MethodPut, "/movie/63e684f9afd8b30f", nil)

	// act

	res, _ := app.Test(req)

	// Assert

	assert.Equal(fiber.StatusBadRequest, res.StatusCode)
}

func TestMarkedAsWatched_Success(t *testing.T) {
	// arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Put("/movie/:id", mark_as_watched_handler.MarkAsWatched(c.GetNetflixRepository()))
	req := httptest.NewRequest(http.MethodPut, "/movie/63e22d6a22d15a2f27e6ba70", nil)

	// act
	res, err := app.Test(req)

	// assert
	assert.NoError(err)
	assert.Equal(fiber.StatusNoContent, res.StatusCode)
}

func TestLoginHandlerInvalidEntity(t *testing.T) {
	// arrange

	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	c := container.New("testing")
	app := fiber.New()
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New()))
	body := []byte(`{
		aaaaa
	}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))

	// act
	res, err := app.Test(req)

	// assert
	assert.NoError(err)
	assert.Equal(fiber.StatusUnprocessableEntity, res.StatusCode)
}

func TestLoginHandlerUsernameNotFound(t *testing.T) {
	// arrange
	data := struct {
		Username string
		Password string
	}{
		Username: "admin",
		Password: "admin123",
	}
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	c := container.New("testing")
	app := fiber.New()
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New()))
	b, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	col := c.GetMongoCollection("users")
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})

	// act
	res, err := app.Test(req)

	// assert
	assert.NoError(err)
	assert.Equal(fiber.StatusBadRequest, res.StatusCode)
}

func TestLoginHandlerInvalidCredentials(t *testing.T) {
	// arrange
	data := struct {
		Username string
		Password string
	}{
		Username: "admin",
		Password: "admin1234",
	}
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	c := container.New("testing")
	app := fiber.New()
	userRepo := c.GetUserRpository()
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New()))
	b, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	col := c.GetMongoCollection("users")
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})
	hashPass, _ := c.GetHashRepository().Hash("admin123")
	_, _ = userRepo.Register(context.Background(), model.User{
		Username:  "admin",
		Password:  string(hashPass),
		Firstname: "Admin",
	})

	// act
	res, err := app.Test(req)

	// assert
	assert.NoError(err)
	assert.Equal(fiber.StatusBadRequest, res.StatusCode)
}

func TestLoginHandlerSuccess(t *testing.T) {
	// arrange
	data := struct {
		Username string
		Password string
	}{
		Username: "admin",
		Password: "admin123",
	}
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	c := container.New("testing")
	app := fiber.New()
	userRepo := c.GetUserRpository()
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New()))
	b, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	col := c.GetMongoCollection("users")
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})

	hashPass, _ := c.GetHashRepository().Hash("admin123")
	id, _ := userRepo.Register(context.Background(), model.User{
		Username:  "admin",
		Password:  string(hashPass),
		Firstname: "Admin",
	})

	// act
	res, err := app.Test(req)
	user1 := model.User{}
	bytes, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	_ = json.Unmarshal(bytes, &user1)

	// assert
	assert.NoError(err)
	assert.Equal(fiber.StatusOK, res.StatusCode)
	assert.Equal(data.Username, user1.Username)
	assert.Equal(id.Hex(), user1.ID.Hex())
}

func TestLoginHandlerValidationFail(t *testing.T) {
	// arrange
	data := struct {
		Username string
		Password string
	}{
		Password: "admin1",
	}
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	c := container.New("testing")
	app := fiber.New()
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New()))
	b, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	// act
	res, err := app.Test(req)

	// assert
	assert.NoError(err)
	assert.Equal(fiber.StatusUnprocessableEntity, res.StatusCode)
}

func TestRegisterHandlerInvalidEntity(t *testing.T) {
	// arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	c := container.New("testing")
	app := fiber.New()
	app.Post("/register", register_handler.Handler(c.GetUserRpository(), c.GetHashRepository(), validator.New()))
	body := []byte(`{
		aaaaa
	}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))

	// act
	res, err := app.Test(req)

	// assert
	assert.NoError(err)
	assert.Equal(fiber.StatusUnprocessableEntity, res.StatusCode)
}

func TestRegisterHandlerSuccess(t *testing.T) {
	// arrange
	data := struct {
		Firstname string
		Username  string
		Password  string
	}{
		Firstname: "Admin",
		Username:  "admin",
		Password:  "admin123",
	}
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	c := container.New("testing")
	app := fiber.New()
	app.Post("/register", register_handler.Handler(c.GetUserRpository(), c.GetHashRepository(), validator.New()))
	b, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(b))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	col := c.GetMongoCollection("users")
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})

	// act
	res, err := app.Test(req)
	user := model.User{}
	bytes, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	_ = json.Unmarshal(bytes, &user)

	// assert
	assert.NoError(err)
	assert.Equal(fiber.StatusCreated, res.StatusCode)
}

func TestRegisterHandlerValidationFail(t *testing.T) {
	// arrange
	data := struct {
		Firstname string
		Username  string
		Password  string
	}{
		Firstname: "Admin",
		Password:  "admin",
	}
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	c := container.New("testing")
	app := fiber.New()
	app.Post("/register", register_handler.Handler(c.GetUserRpository(), c.GetHashRepository(), validator.New()))
	b, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(b))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	// act
	res, err := app.Test(req)

	// assert
	assert.NoError(err)
	assert.Equal(fiber.StatusUnprocessableEntity, res.StatusCode)
}
