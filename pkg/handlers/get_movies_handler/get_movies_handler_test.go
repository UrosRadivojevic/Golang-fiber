package get_movies_handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/handlers/get_movies_handler"
	"github.com/urosradivojevic/health/pkg/handlers/login_handler.go"
	"github.com/urosradivojevic/health/pkg/middleware/jwt_middleware"
	"github.com/urosradivojevic/health/testing_utils/helpers"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGetMovies_Success(t *testing.T) {
	// arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Get("/movies", jwt_middleware.AuthRequired(), get_movies_handler.GetMovies(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New(), c.GetTokenService()))
	col := c.GetMongoCollection("users")
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})
	token := helpers.RegisterAndLoginUser(c.GetUserRpository(), c.GetHashRepository(), t, app)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
	// act
	res, err := app.Test(req)

	// assert
	assert.NoError(err)
	assert.Equal(fiber.StatusOK, res.StatusCode)
}
