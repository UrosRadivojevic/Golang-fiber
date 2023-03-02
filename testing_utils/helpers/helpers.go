package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/handlers/login_handler.go"
	"github.com/urosradivojevic/health/pkg/model"
	"github.com/urosradivojevic/health/pkg/repositories/user_repository"
	"github.com/urosradivojevic/health/pkg/services/hasher"
)

func RegisterAndLoginUser(userRepo user_repository.Interface, hasher hasher.Interface, t *testing.T, app *fiber.App) string {
	t.Helper()
	hashPass, _ := hasher.Hash("admin123")
	_, _ = userRepo.Register(context.Background(), model.User{
		Username:  "admin",
		Password:  string(hashPass),
		Firstname: "Admin",
	})
	b := []byte(`{"username": "admin", "password":"admin123"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(b))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	res, _ := app.Test(req)
	bytes, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	var m login_handler.Response
	_ = json.Unmarshal(bytes, &m)
	return m.Token
}
