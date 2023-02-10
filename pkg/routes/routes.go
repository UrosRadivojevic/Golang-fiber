package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/handler"
)

func SetUpRoutes(app *fiber.App, c *container.Container) {
	app.Delete("/movie/:id", handler.DeleteMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	app.Get("/movie/:id", handler.GetMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	// app.Get("/health", handler.Health2(time.Now()))
	app.Put("/movie/:id", handler.MarkAsWatched(c.GetNetflixRepository()))
	app.Post("/movie", handler.CreateMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	app.Get("/movies", handler.GetMovies(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
}
