package get_movie_handler_test

import (
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
	"github.com/urosradivojevic/health/pkg/handlers/get_movie_handler"
	"github.com/urosradivojevic/health/pkg/handlers/login_handler.go"
	"github.com/urosradivojevic/health/pkg/middleware/jwt_middleware"
	"github.com/urosradivojevic/health/pkg/model"
	"github.com/urosradivojevic/health/testing_utils/helpers"
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
		ID:       "63ff1f4a34a7bdac8e42958d",
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
	app.Get("/movie/:id", jwt_middleware.AuthRequired(), get_movie_handler.GetMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New(), c.GetTokenService()))
	req := httptest.NewRequest(http.MethodGet, "/movie/63ff1f4a34a7bdac8e42958d", nil)
	col := c.GetMongoCollection("users")
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})
	token := helpers.RegisterAndLoginUser(c.GetUserRpository(), c.GetHashRepository(), t, app)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)

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
	app.Get("/movie/:id", jwt_middleware.AuthRequired(), get_movie_handler.GetMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New(), c.GetTokenService()))
	req := httptest.NewRequest(http.MethodGet, "/movie/63e684f9afd8b30f56511", nil)
	col := c.GetMongoCollection("users")
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})
	token := helpers.RegisterAndLoginUser(c.GetUserRpository(), c.GetHashRepository(), t, app)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
	// act

	res, _ := app.Test(req)

	// Assert

	assert.Equal(fiber.StatusBadRequest, res.StatusCode)
}
