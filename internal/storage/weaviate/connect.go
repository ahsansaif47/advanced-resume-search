package weaviate

import (
	"context"
	"log"
	"sync"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	// "github.com/weaviate/weaviate-go-client/v5/weaviate/schema"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

type Database struct {
	Client *weaviate.Client
}

var dbInstance *Database
var once sync.Once

func GetDatabaseConnection() *Database {
	once.Do(func() {
		dbInstance = &Database{
			Client: ConnectWeaviate(),
		}
	})
	return dbInstance
}

func ConnectWeaviate() *weaviate.Client {
	cfg := weaviate.Config{
		Host:   "localhost:8081",
		Scheme: "http",
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		log.Fatalf("error creating weaviate client: %v", err)
	}

	live, err := client.Misc().LiveChecker().Do(context.Background())
	if err != nil {
		log.Fatalf("error checking live status of weaviate: %v", err)
	}

	if !live {
		log.Fatal("Error connecting to weaviate!")
	}

	return client
}

func CreateSchema(c *weaviate.Client, className string) {

	resumeClass := &models.Class{
		Class:      className,
		Vectorizer: "text2vec-transformers",
		ModuleConfig: map[string]any{
			"text2vec-transformers": map[string]any{
				"pooling":            "masked_mean",
				"vectorizeClassName": false,
			},
		},
		Properties: []*models.Property{
			{
				Name:     "personal_information",
				DataType: schema.DataTypeObject.PropString(),
				// NestedProperties: []*models.NestedProperty{},
			},
			{
				Name:     "education",
				DataType: schema.DataTypeObjectArray.PropString(),
			},
			{
				Name:     "work_experience",
				DataType: schema.DataTypeObjectArray.PropString(),
			},
			{
				Name:     "skills",
				DataType: schema.DataTypeStringArray.PropString(),
			},
			{
				Name:     "extra",
				DataType: schema.DataTypeObject.PropString(),
			},
		},
		VectorIndexType: "hnsw",
	}

	if exists, err := c.Schema().ClassExistenceChecker().WithClassName(className).Do(context.Background()); err != nil && !exists {
		err := c.Schema().ClassCreator().WithClass(resumeClass).Do(context.Background())
		if err != nil {
			log.Fatalf("error creating class: %v", err.Error())
		}
	}
}
