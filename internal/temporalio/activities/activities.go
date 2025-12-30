package activities

import (
	"context"

	"github.com/ahsansaif47/advanced-resume/integrations/gemini"
	repo "github.com/ahsansaif47/advanced-resume/internal/storage/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate"
)

type Activities struct {
	WeaviateClient *weaviate.Client
	GenAIClient    gemini.IGeminiClient
}

func NewActivities(ctx context.Context) (*Activities, error) {
	client, err := repo.ConnectWeaviate()
	if err != nil {
		return nil, err
	}
	geminiClient, err := gemini.NewGeminiClient()

	if err != nil {
		return nil, err
	}

	return &Activities{
		WeaviateClient: client,
		GenAIClient:    geminiClient,
	}, nil
}
