package get_movies_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/cache"
	"github.com/urosradivojevic/health/pkg/repositories"
)

// ShowAccount godoc
//
//		@Summary		  Get movies
//		@Description	Returns all movies from db
//		@Tags			    movies
//		@Accept			  json
//		@Produce		  json
//		@Success		 200	{object}	[]model.Netflix
//	@Router			 /movies [get]
func GetMovies(repo repositories.NetflixInterface, redis cache.RedisCacheInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		movies, err := repo.GetAllMovies()

		if err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(movies)
	}
}
