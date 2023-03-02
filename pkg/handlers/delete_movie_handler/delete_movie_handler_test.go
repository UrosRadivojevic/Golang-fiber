package delete_movie_handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/handlers/delete_movie_handler"
	"github.com/urosradivojevic/health/pkg/handlers/login_handler.go"
	"github.com/urosradivojevic/health/pkg/middleware/jwt_middleware"
	"github.com/urosradivojevic/health/testing_utils/helpers"
	"go.mongodb.org/mongo-driver/bson"
)

func TestDeleteMovie_Success(t *testing.T) {
	// arrange
	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Delete("/movie/:id", jwt_middleware.AuthRequired(), delete_movie_handler.DeleteMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	req := httptest.NewRequest(http.MethodDelete, "/movie/63e684f9afd8b30f56511e46", nil)
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
	app.Delete("/movie/:id", jwt_middleware.AuthRequired(), delete_movie_handler.DeleteMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	req := httptest.NewRequest(http.MethodDelete, "/movie/63e684f9afd8b30f", nil)
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New(), c.GetTokenService()))
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
