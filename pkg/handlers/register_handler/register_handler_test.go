package register_handler_test

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
	"github.com/urosradivojevic/health/pkg/handlers/register_handler"
	"github.com/urosradivojevic/health/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

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
