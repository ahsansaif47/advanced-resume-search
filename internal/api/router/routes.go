package router

import (
	"context"
	"log"

	"github.com/ahsansaif47/advanced-resume/integrations/gemini"
	"github.com/ahsansaif47/advanced-resume/internal/api/controllers"
	"github.com/ahsansaif47/advanced-resume/internal/api/handlers"
	weviateRepo "github.com/ahsansaif47/advanced-resume/internal/storage/weaviate"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"go.temporal.io/sdk/client"
)

// @title						Talent Pooling Local API
// @version					1.0
// @description				This is a swagger for Talent Pooling
// @host						localhost:9090
// @BasePath					/api/v1
// @schemes					http
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func InitRoutes(app *fiber.App, db *weaviate.Client, tempClient client.Client) {
	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	aiClient, err := gemini.NewGeminiClient()
	if err != nil {
		log.Fatalf("Error Initializing Routes: %s", err.Error())
	}

	repo := weviateRepo.NewWeviateRepository(context.Background(), db)
	service := controllers.NewWeaviateService(repo, aiClient, tempClient)
	handlers := handlers.NewHandlers(service)

	v1.Post("/upload", handlers.UploadResume)
	v1.Get("/search", handlers.VectorSearch)
}
