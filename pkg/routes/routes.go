package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/handlers/create_movie_handler"
	"github.com/urosradivojevic/health/pkg/handlers/delete_movie_handler"
	"github.com/urosradivojevic/health/pkg/handlers/get_movie_handler"
	"github.com/urosradivojevic/health/pkg/handlers/get_movies_handler"
	"github.com/urosradivojevic/health/pkg/handlers/mark_as_watched_handler"
)

func SetUpRoutes(app *fiber.App, c *container.Container) {
	app.Delete("/movie/:id", delete_movie_handler.DeleteMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	app.Get("/movie/:id", get_movie_handler.GetMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	// app.Get("/health", handler.Health2(time.Now()))
	app.Put("/movie/:id", mark_as_watched_handler.MarkAsWatched(c.GetNetflixRepository()))
	app.Post("/movie", create_movie_handler.CreateMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	app.Get("/movies", get_movies_handler.GetMovies(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
}
