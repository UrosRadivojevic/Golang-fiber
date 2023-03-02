package routes

import (
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/handlers/create_movie_handler"
	"github.com/urosradivojevic/health/pkg/handlers/delete_movie_handler"
	"github.com/urosradivojevic/health/pkg/handlers/get_movie_handler"
	"github.com/urosradivojevic/health/pkg/handlers/get_movies_handler"
	"github.com/urosradivojevic/health/pkg/handlers/login_handler.go"
	"github.com/urosradivojevic/health/pkg/handlers/mark_as_watched_handler"
	"github.com/urosradivojevic/health/pkg/handlers/register_handler"
	"github.com/urosradivojevic/health/pkg/middleware/jwt_middleware"
)

func SetUpRoutes(app *fiber.App, c *container.Container) {
	validate := validator.New()
	app.Delete("/movie/:id", jwt_middleware.AuthRequired(), delete_movie_handler.DeleteMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	app.Get("/movie/:id", jwt_middleware.AuthRequired(), get_movie_handler.GetMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	app.Put("/movie/:id", jwt_middleware.AuthRequired(), mark_as_watched_handler.MarkAsWatched(c.GetNetflixRepository()))
	app.Post("/movie", jwt_middleware.AuthRequired(), create_movie_handler.CreateMovie(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	app.Get("/movies", jwt_middleware.AuthRequired(), get_movies_handler.GetMovies(c.GetNetflixRepository(), c.GetRedisCacheRepository()))
	app.Post("/register", register_handler.Handler(c.GetUserRpository(), c.GetHashRepository(), validate))
	app.Post("/login", login_handler.Handler(c.GetLoginRepository(), validate, c.GetTokenService()))
}
