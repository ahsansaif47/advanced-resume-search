package weaviate

import (
	"context"
	"log"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

type IWeaviateRepository interface {
	AddResumeToDB(className string, props map[string]any) (string, error)
	BatchUploadResume(className string, items []map[string]any) error
	VectorSearch(className, query string) (map[string]models.JSONObject, error)
}

type WeaviateRepository struct {
	ctx    context.Context
	Client *weaviate.Client
}

func NewWeviateRepository(ctx context.Context, client *weaviate.Client) IWeaviateRepository {
	return &WeaviateRepository{
		ctx:    ctx,
		Client: client,
	}
}

func (r *WeaviateRepository) AddResumeToDB(className string, props map[string]any) (string, error) {

	resp, err := r.Client.Data().Creator().WithClassName(className).WithProperties(props).Do(r.ctx)
	if err != nil {
		return "", err
	}

	return resp.Object.ID.String(), nil
}

func (r *WeaviateRepository) BatchUploadResume(className string, items []map[string]any) error {

	batch := r.Client.Batch().ObjectsBatcher()
	for _, item := range items {
		batch = batch.WithObjects(&models.Object{
			Class:      className,
			Properties: item,
		})
	}

	_, err := batch.Do(r.ctx)
	log.Printf("Imported & vectorized %d objects into the %s collection\n", len(items), className)
	return err
}

func (r *WeaviateRepository) VectorSearch(className, query string) (map[string]models.JSONObject, error) {
	response, err := r.Client.GraphQL().Get().
		WithClassName(className).
		WithNearText(r.Client.GraphQL().NearTextArgBuilder().WithConcepts([]string{query})).
		WithFields(
			graphql.Field{
				Name: "personal_information { name title email phone github linkedin }",
			},
			graphql.Field{
				Name: "education { degree institution dates location }",
			},
			graphql.Field{
				Name: "work_experience { company title dates responsibilities }",
			},
			graphql.Field{
				Name: "skills",
			},
			graphql.Field{
				Name: "_additional { id score distance }",
			},
		).
		WithLimit(20).
		Do(context.Background())

	return response.Data, err

}
