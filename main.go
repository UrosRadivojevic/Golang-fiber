package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gkampitakis/fiber-modules/gracefulshutdown"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/urosradivojevic/health/container"
	"github.com/urosradivojevic/health/handler"
)

var flagvar string

func init() {
	flag.StringVar(&flagvar, "env", "development", "help message for flagname")
	flag.Parse()
}

func main() {

	c := container.New(flagvar)
	fmt.Printf("%v\n", flagvar)
	app := fiber.New()
	env := c.GetEnviorment()
	if env == "development" {
		err := godotenv.Load(".env.development")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	if env == "testing" {
		err := godotenv.Load(".env.testing")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	if env == "production" {
		fmt.Println("No access")
	}
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Div Rhino!")
	})

	app.Delete("/movie/:id", handler.DeleteMovie(c.GetNetflixRepository()))
	app.Get("/health", handler.Health2(time.Now()))
	app.Put("/movie/:id", handler.MarkAsWatched(c.GetNetflixRepository()))
	app.Post("/movie", handler.CreateMovie(c.GetNetflixRepository()))
	app.Get("/movies", handler.GetMovies(c.GetNetflixRepository()))

	fmt.Println("Server is getting started...")
	gracefulshutdown.Listen(app, "localhost:3000", gracefulshutdown.Default())
	// GracefullShoutdown(app)
	// err := app.Listen(":3000")
	// if err != nil {
	// 	fmt.Println(err)
	// }

}

// func GracefullShoutdown(app *fiber.App) {
// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt)

// 	serverShutdown := make(chan struct{})

// 	go func() {
// 		_ = <-c
// 		fmt.Println("Gracefully shutting down...")
// 		_ = app.Shutdown()
// 		serverShutdown <- struct{}{}
// 	}()

// 	if err := app.Listen(":3000"); err != nil {
// 		log.Panic(err)
// 	}

// 	<-serverShutdown

// 	fmt.Println("Running cleanup tasks...")
// 	// Your cleanup tasks go here
// }