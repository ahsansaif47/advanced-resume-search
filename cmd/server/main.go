package main

import (
	"fmt"
	"log"

	"github.com/ahsansaif47/advanced-resume/config"
	"github.com/ahsansaif47/advanced-resume/internal/api/router"
	"github.com/ahsansaif47/advanced-resume/internal/storage/weaviate"
	"github.com/ahsansaif47/advanced-resume/internal/temporalio"
	"github.com/ahsansaif47/advanced-resume/internal/temporalio/workflows"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.temporal.io/sdk/client"

	_ "github.com/ahsansaif47/advanced-resume/internal/api/docs"
)

// NOTE: Start temporal with the command below before running tests...
// temporal server start-dev --ui-port 8080

var className = "Resume"

func main() {

	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	weaviateClient := weaviate.ConnectWeaviate()
	weaviate.CreateSchema(weaviateClient, className)

	tempClient, errCh := temporalio.StartWorker()
	go func(client client.Client) {
		for err := range errCh {
			log.Printf("Worker error: %v", err)

			// We need to use retry here to fail incase there is some connectivity error.
			// if errors.Is(err, <connection-error>), then log.Fatal...
			client.Close()
			client, errCh = temporalio.StartWorker()
		}
	}(tempClient)

	// print(<-errCh) // doing this temporarily... upper logic is the correct one though incomplete
	defer (tempClient).Close()

	router.InitRoutes(app, weaviateClient, tempClient)
	port := config.GetConfig().Port
	workflows.ExecuteWorkflow_StoreResumeToWeaviate(tempClient, "/home/ahsansaif/Downloads/AhsanResume202507.pdf")
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))

}
