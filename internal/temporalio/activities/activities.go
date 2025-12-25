package activities

import (
	"context"

	repo "github.com/ahsansaif47/advanced-resume/internal/storage/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

type Activities struct {
	WeaviateClient *weaviate.Client
}

func NewActivities(ctx context.Context) (*Activities, error) {
	client, err := repo.ConnectWeaviate()
	if err != nil {
		return nil, err
	}

	return &Activities{
		WeaviateClient: client,
	}, nil
}
