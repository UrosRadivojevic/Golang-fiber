package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/handler"
)

func SetUpRoutes(app *fiber.App, c *container.Container) {
	app.Delete("/movie/:id", handler.DeleteMovie(c.GetNetflixRepository()))
	app.Get("/health", handler.Health2(time.Now()))
	app.Put("/movie/:id", handler.MarkAsWatched(c.GetNetflixRepository()))
	app.Post("/movie", handler.CreateMovie(c.GetNetflixRepository()))
	app.Get("/movies", handler.GetMovies(c.GetNetflixRepository()))
}
