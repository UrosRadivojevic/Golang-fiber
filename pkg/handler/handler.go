package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/pkg/cache"
	"github.com/urosradivojevic/health/pkg/model"
	"github.com/urosradivojevic/health/pkg/repositories"
)

// var startTime time.Time

// func Health2(stratime time.Time) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		health := true
// 		c.SendString("Hello, World!\n")
// 		uptime := time.Since(startTime)
// 		return c.JSON(fiber.Map{
// 			"healthy": health,
// 			"update":  uptime.Seconds(),
// 		})
// 	}
// }

func CreateMovie(repo repositories.NetflixInterface, redis cache.RedisCacheInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var movie model.Netflix
		if err := c.BodyParser(&movie); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"errors": err.Error(),
			})
		}
		movieId, err := repo.InsertOneMovie(movie)
		if err != nil {
			return err
		}
		movie.ID = movieId
		err = redis.SetMovie(c.UserContext(), movie)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(movie)
	}
}

func GetMovie(repo repositories.NetflixInterface, redis cache.RedisCacheInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		movieId := c.Params("id")
		if len(movieId) != 24 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid mongoID"})
		}
		m := model.Netflix{}
		var movie1 model.Netflix
		movie, err := redis.Get(c.UserContext(), movieId)
		if err != nil {
			fmt.Println("Object not in cache, searcing in database. Error:", err)
		}
		if movie == m {
			movie1, err = repo.GetOneMovie(movieId)
			if err != nil {
				return err
			}
			if err = redis.SetMovie(c.UserContext(), movie1); err != nil {
				return err
			}

			fmt.Println("Object returned from database and inserted in cache.")
			return c.JSON(movie1)
		}
		fmt.Println("Object returned from cache.")
		return c.JSON(movie)
	}
}

func GetMovies(repo repositories.NetflixInterface, redis cache.RedisCacheInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		movies, err := repo.GetAllMovies()

		if err != nil {
			return err
		}
		file, _ := json.MarshalIndent(movies, "", " ")
		_ = ioutil.WriteFile("test.json", file, 0o644)
		return c.Status(fiber.StatusOK).JSON(movies)
	}
}

func MarkAsWatched(repo repositories.NetflixInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		movieId := c.Params("id")
		if len(movieId) != 24 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid mongoID"})
		}
		err := repo.UpdateOneMovie(movieId)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusNoContent).JSON(movieId)
	}
}

func DeleteMovie(repo repositories.NetflixInterface, redis cache.RedisCacheInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		movieId := c.Params("id")
		if len(movieId) != 24 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "Invalid mongoID"})
		}
		if err := redis.Delete(c.UserContext(), movieId); err != nil {
			return err
		}
		if err := repo.DeleteOneMovie(movieId); err != nil {
			return err
		}
		return c.Status(fiber.StatusNoContent).JSON(movieId)
	}
}
