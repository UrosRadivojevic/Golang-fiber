package delete_movie_handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/cache"
	"github.com/urosradivojevic/health/pkg/message"
	"github.com/urosradivojevic/health/pkg/rabbitmq/publisher"
	"github.com/urosradivojevic/health/pkg/repositories/movie_repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShowAccount godoc
//
//			@Summary		  Delete movie
//			@Description	Deletes movie from db by movieID
//			@Tags			    movies
//			@Accept			  json
//			@Produce		  json
//			@Success		 204
//		 @Failure      	 400   {object}    message.Msg "Invalid Object ID"
//	   @Param id   path string true "Movie ID"
//		@Router			 /movie/{id} [delete]
func DeleteMovie(repo movie_repository.NetflixInterface, redis cache.RedisCacheInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		movieId := c.Params("id")
		if !primitive.IsValidObjectID(movieId) {
			return c.Status(fiber.StatusBadRequest).JSON(message.Msg{
				Message: "Invalid ID.",
			})
		}
		if err := redis.Delete(c.UserContext(), movieId); err != nil {
			return err
		}
		if err := repo.DeleteOneMovie(movieId); err != nil {
			return err
		}
		publisher.Publish(movieId, "campaings_deleted_queue")
		return c.SendStatus(fiber.StatusNoContent)
	}
}
