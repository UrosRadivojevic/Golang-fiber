package mark_as_watched_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/repositories/movie_repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type message struct {
	Message string `json:"message"`
}

// ShowAccount godoc
//
//			@Summary		  Update movie
//			@Description	Updates movie filed watched to true
//			@Tags			    movies
//			@Accept			  json
//			@Produce		  json
//			@Success		 204
//		 @Failure      	 400   {object}    message "Invalid Object ID"
//	   @Param id   path string true "Movie ID"
//		@Router			 /movie/{id} [put]
func MarkAsWatched(repo movie_repository.NetflixInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		movieId := c.Params("id")
		if !primitive.IsValidObjectID(movieId) {
			return c.Status(fiber.StatusBadRequest).JSON(message{
				Message: "Invalid ID.",
			})
		}
		err := repo.UpdateOneMovie(movieId)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
