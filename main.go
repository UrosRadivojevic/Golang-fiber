package main

import (
	"flag"
	"fmt"
	"log"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gkampitakis/fiber-modules/gracefulshutdown"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/urosradivojevic/health/docs"
	"github.com/urosradivojevic/health/pkg/container"
	"github.com/urosradivojevic/health/pkg/routes"
)

var flagvar string

func init() {
	flag.StringVar(&flagvar, "env", "development", "help message for flagname")
	flag.Parse()
}

// @title			Fiber Example API
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
	app.Get("/swagger/*", swagger.HandlerDefault)

	routes.SetUpRoutes(app, c)

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
// kanal sluzi da sinhronizujem dve niti
// go func, sponujem novu nit
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
