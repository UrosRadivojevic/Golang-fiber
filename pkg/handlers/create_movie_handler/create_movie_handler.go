package create_movie_handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/cache"
	"github.com/urosradivojevic/health/pkg/model"
	"github.com/urosradivojevic/health/pkg/repositories"
	"github.com/urosradivojevic/health/pkg/requests"
)

type message struct {
	Message string `json:"message"`
}

// ShowAccount godoc
//
//		@Summary		  Create movie
//		@Description	Create movie
//		@Tags			    movies
//		@Accept			  json
//		@Produce		  json
//		@Success		 201	{object}	model.Netflix
//	 @Failure      	 422   {object}    message "Validation failed"
//	 @Param request body requests.CreateMovieRequest true "Movie"
//		@Router			 /movie [post]
func CreateMovie(repo repositories.NetflixInterface, redis cache.RedisCacheInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request requests.CreateMovieRequest

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(message{
				Message: err.Error(),
			})
		}
		validate := validator.New()
		err := validate.Struct(request)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(message{
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

		return c.Status(fiber.StatusCreated).JSON(movie)
	}
}
