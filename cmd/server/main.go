package main

import (
	"fmt"
	"log"

	"github.com/ahsansaif47/advanced-resume/config"
	"github.com/ahsansaif47/advanced-resume/internal/api/router"
	"github.com/ahsansaif47/advanced-resume/internal/storage/weaviate"
	"github.com/ahsansaif47/advanced-resume/internal/temporalio"
	"github.com/gofiber/fiber/v2"
	"go.temporal.io/sdk/client"
)

// NOTE: Start temporal with the command below before running tests...
// temporal server start-dev --ui-port 8080

var className = "Resume"

func main() {

	app := fiber.New()
	weaviateClient := weaviate.ConnectWeaviate()
	weaviate.CreateSchema(weaviateClient, className)

	tempClient, errCh := temporalio.StartWorker()
	go func(client *client.Client) {
		for err := range errCh {
			log.Printf("Worker error: %v", err)

			// We need to use retry here to fail incase there is some connectivity error.
			// if errors.Is(err, <connection-error>), then log.Fatal...
			(*client).Close()
			client, errCh = temporalio.StartWorker()
		}
	}(tempClient)

	print(<-errCh) // doing this temporarily... upper logic is the correct one though incomplete
	defer (*tempClient).Close()

	router.InitRoutes(app, weaviateClient, tempClient)
	port := config.GetConfig().Port
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
