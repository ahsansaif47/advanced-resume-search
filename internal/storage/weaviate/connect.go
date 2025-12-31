package weaviate

import (
	"context"
	"fmt"
	"sync"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

type Database struct {
	Client *weaviate.Client
}

var dbInstance *Database
var once sync.Once

func GetDatabaseConnection() *Database {
	once.Do(func() {
		client, _ := ConnectWeaviate()
		dbInstance = &Database{
			Client: client,
		}
	})
	return dbInstance
}

func ConnectWeaviate() (*weaviate.Client, error) {
	cfg := weaviate.Config{
		Host:   "localhost:8081",
		Scheme: "http",
	}

	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating weaviate client: %s", err.Error())
		// log.Fatalf("error creating weaviate client: %v", err)
	}

	live, err := client.Misc().LiveChecker().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("checking weaviate live status: %s", err.Error())
		// log.Fatalf("error checking live status of weaviate: %v", err)
	}

	if !live {
		return nil, fmt.Errorf("weaviate is not live")
		// log.Fatal("Error connecting to weaviate!")
	}

	return client, nil
}
