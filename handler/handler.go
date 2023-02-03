package handler

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/urosradivojevic/health/model"
	"github.com/urosradivojevic/health/repositories"
)

var startTime time.Time

func Health2(stratime time.Time) fiber.Handler {
	return func(c *fiber.Ctx) error {
		health := true
		c.SendString("Hello, World!\n")
		uptime := time.Since(startTime)
		return c.JSON(fiber.Map{
			"healthy": health,
			"update":  uptime.Seconds(),
		})
	}
}
func CreateMovie(repo repositories.NetflixInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var movie model.Netflix
		if err := c.BodyParser(&movie); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"errors": err.Error(),
			})
		}
		err := repo.InsertOneMovie(movie)
		if err != nil {
			return err
		}
		return c.JSON(movie)
	}
}

func GetMovies(repo repositories.NetflixInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		movies, err := repo.GetAllMovies()
		if err != nil {
			return err
		}
		file, _ := json.MarshalIndent(movies, "", " ")
		_ = ioutil.WriteFile("test.json", file, 0644)
		return c.JSON(movies)
	}
}

func MarkAsWatched(repo repositories.NetflixInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		movieId := c.Params("id")
		err := repo.UpdateOneMovie(movieId)
		if err != nil {
			return err
		}

		return c.JSON(movieId)
	}
}

func DeleteMovie(repo repositories.NetflixInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		movieId := c.Params("id")
		err := repo.DeleteOneMovie(movieId)
		if err != nil {
			return err
		}
		return c.JSON(movieId)
	}
}
