package login_handler_test

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
	"github.com/urosradivojevic/health/pkg/handlers/login_handler.go"
	"github.com/urosradivojevic/health/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

func TestLoginHandlerInvalidEntity(t *testing.T) {
	// arrange

	assert := require.New(t)
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	c := container.New("testing")
	app := fiber.New()
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New(), c.GetTokenService()))
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
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New(), c.GetTokenService()))
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
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New(), c.GetTokenService()))
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
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New(), c.GetTokenService()))
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
	var m login_handler.Response

	// act
	res, err := app.Test(req)
	user1 := model.User{}
	bytes, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	_ = json.Unmarshal(bytes, &m)
	user1 = m.User

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
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validator.New(), c.GetTokenService()))
	b, _ := json.Marshal(data)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

	// act
	res, err := app.Test(req)

	// assert
	assert.NoError(err)
	assert.Equal(fiber.StatusUnprocessableEntity, res.StatusCode)
}
