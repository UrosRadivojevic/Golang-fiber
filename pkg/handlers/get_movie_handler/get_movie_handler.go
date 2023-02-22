package get_movie_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/urosradivojevic/health/pkg/cache"
	"github.com/urosradivojevic/health/pkg/message"
	"github.com/urosradivojevic/health/pkg/model"
	"github.com/urosradivojevic/health/pkg/repositories/movie_repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShowAccount godoc
//
//			@Summary		  Get movie
//			@Description	Returns movie from db by movieID
//			@Tags			    movies
//			@Accept			  json
//			@Produce		  json
//			@Success		 200	{object}	model.Netflix
//		 @Failure      	 400   {object}    message.Msg "Invalid object ID"
//	   @Param id   path string true "Movie ID" minlength(24) maxlength(24)
//		@Router			 /movie/{id} [get]
func GetMovie(repo movie_repository.NetflixInterface, redis cache.RedisCacheInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		movieId := c.Params("id")
		if !primitive.IsValidObjectID(movieId) {
			return c.Status(fiber.StatusBadRequest).JSON(message.Msg{
				Message: "Invalid ID.",
			})
		}
		m := model.Netflix{}
		var movie1 model.Netflix
		movie, err := redis.Get(c.UserContext(), movieId)
		if err != nil {
			log.Error().Err(err).Msg("Object not in cache, searcing in database.")
		}
		if movie == m {
			movie1, err = repo.GetOneMovie(movieId)
			if err != nil {
				return err
			}
			if err = redis.SetMovie(c.UserContext(), movie1); err != nil {
				return err
			}

			log.Info().Msg("Object returned from database.")
			return c.Status(fiber.StatusOK).JSON(movie1)
		}
		log.Info().Msg("Object return from cache.")
		return c.Status(fiber.StatusOK).JSON(movie)
	}
}
