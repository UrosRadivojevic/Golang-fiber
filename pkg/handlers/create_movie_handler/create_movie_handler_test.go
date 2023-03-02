package create_movie_handler_test

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
	"github.com/urosradivojevic/health/pkg/handlers/login_handler.go"
	"github.com/urosradivojevic/health/pkg/middleware/jwt_middleware"
	"github.com/urosradivojevic/health/pkg/model"
	"github.com/urosradivojevic/health/testing_utils/helpers"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateMovie_InvalidEntity(t *testing.T) {
	// arrange

	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	c := container.New("testing")
	app := fiber.New()
	app.Post("/movie", jwt_middleware.AuthRequired(), create_movie_handler.CreateMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))

	body := []byte(`{
		aaaaa
	}`)
	req := httptest.NewRequest(http.MethodPost, "/movie", bytes.NewBuffer(body))
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New(), c.GetTokenService()))
	col := c.GetMongoCollection("users")
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})
	token := helpers.RegisterAndLoginUser(c.GetUserRpository(), c.GetHashRepository(), t, app)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
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
	app.Post("/movie", jwt_middleware.AuthRequired(), create_movie_handler.CreateMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))

	b, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/movie", bytes.NewBuffer(b))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New(), c.GetTokenService()))
	col := c.GetMongoCollection("users")
	t.Cleanup(func() {
		_, _ = col.DeleteMany(context.Background(), bson.D{})
	})
	token := helpers.RegisterAndLoginUser(c.GetUserRpository(), c.GetHashRepository(), t, app)
	req.Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
	// Act
	res, err := app.Test(req)

	movie := model.Netflix{}
	bytes, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	_ = json.Unmarshal(bytes, &movie)

	// Assert

	assert.NoError(err)
	assert.Equal(fiber.StatusCreated, res.StatusCode)
	assert.Equal("DieHard", movie.Movie)
	assert.Equal(2005, movie.Year)
	assert.Equal(true, movie.Watched)
	assert.Equal("Bruce", movie.LeadRole)
}
