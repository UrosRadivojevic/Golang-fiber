package handler_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/handler"
)

func TestGetMovies_Success(t *testing.T) {
	t.Parallel()
	t.Setenv("MONGODB_URI", "mongodb://localhost:27017/")
	t.Setenv("MONGODB_DB", "netflix")
	c := container.New("testing")
	app := fiber.New()
	app.Get("/movies", handler.GetMovies(c.GetNetflixRepository()))
}
