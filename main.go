package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gkampitakis/fiber-modules/gracefulshutdown"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "github.com/urosradivojevic/health/docs"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/routes"
)

var flagvar string

func init() {
	flag.StringVar(&flagvar, "env", "development", "help message for flagname")
	flag.Parse()
}

// @title			Fiber Movies API
// @version		1.0
// @description	This is a sample swagger for Fiber
// @termsOfService	http://swagger.io/terms/
// @contact.name	API Support
// @contact.email	fiber@swagger.io
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			localhost:3000
// @BasePath		/
func main() {
	c := container.New(flagvar)
	fmt.Printf("%v\n", flagvar)
	app := fiber.New()
	env := c.GetEnviorment()
	if env == "development" {
		err := godotenv.Load(".env.development")
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		if err != nil {
			log.Fatal()
		}
	}
	if env == "testing" {
		err := godotenv.Load(".env.testing")
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		if err != nil {
			log.Fatal()
		}
	}
	if env != "production" {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	routes.SetUpRoutes(app, c)
	portString := os.Getenv("PORT")
	port, err := strconv.Atoi(portString)
	if err != nil {
		log.Fatal().Err(err).Msg("Port is not a number.")
	}
	log.Info().Int("port", port).Msg("Server is getting started.")

	gracefulshutdown.Listen(app, "localhost:3000", gracefulshutdown.Default())
}
