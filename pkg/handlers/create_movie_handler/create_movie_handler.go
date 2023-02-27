package create_movie_handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/cache"
	"github.com/urosradivojevic/health/pkg/message"
	"github.com/urosradivojevic/health/pkg/model"
	"github.com/urosradivojevic/health/pkg/rabbitmq/publisher"
	"github.com/urosradivojevic/health/pkg/repositories/movie_repository"
	"github.com/urosradivojevic/health/pkg/requests"
)

// ShowAccount godoc
//
//		@Summary		  Create movie
//		@Description	Create movie
//		@Tags			    movies
//		@Accept			  json
//		@Produce		  json
//		@Success		 201	{object}	model.Netflix
//	 @Failure      	 422   {object}    message.Msg "Validation failed"
//	 @Param request body requests.CreateMovieRequest true "Movie"
//		@Router			 /movie [post]
func CreateMovie(repo movie_repository.NetflixInterface, redis cache.RedisCacheInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request requests.CreateMovieRequest

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(message.Msg{
				Message: err.Error(),
			})
		}
		validate := validator.New()
		err := validate.Struct(request)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(message.Msg{
				Message: "Validation failed",
			})
		}
		movieId, err := repo.InsertOneMovie(request)
		if err != nil {
			return err
		}
		movie := model.Netflix{
			ID:       movieId,
			Movie:    request.Movie,
			Watched:  request.Watched,
			Year:     request.Year,
			LeadRole: request.LeadRole,
		}

		err = redis.SetMovie(c.UserContext(), movie)
		if err != nil {
			return err
		}
		publisher.Publish(movie.ID.Hex(), "campaings_created_queue")
		return c.Status(fiber.StatusCreated).JSON(movie)
	}
}
