package models

import (
	"context"
	"log"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

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
				Name:     "summary",
				DataType: schema.DataTypeText.PropString(),
			},
			{
				Name:     "personal_information",
				DataType: schema.DataTypeObject.PropString(),
				NestedProperties: []*models.NestedProperty{
					{Name: "name", DataType: schema.DataTypeText.PropString()},
					{Name: "email", DataType: schema.DataTypeText.PropString()},
					{Name: "phone", DataType: schema.DataTypeText.PropString()},
					{Name: "github", DataType: schema.DataTypeText.PropString()},
					{Name: "linkedin", DataType: schema.DataTypeText.PropString()},
				},
			},
			{
				Name:     "education",
				DataType: schema.DataTypeObjectArray.PropString(),
				NestedProperties: []*models.NestedProperty{
					{Name: "institution", DataType: schema.DataTypeText.PropString()},
					{Name: "degree", DataType: schema.DataTypeText.PropString()},
					{Name: "dates", DataType: schema.DataTypeText.PropString()},
				},
			},
			{
				Name:     "work_experience",
				DataType: schema.DataTypeObjectArray.PropString(),
				NestedProperties: []*models.NestedProperty{
					{Name: "company", DataType: schema.DataTypeText.PropString()},
					{Name: "title", DataType: schema.DataTypeText.PropString()},
					{Name: "dates", DataType: schema.DataTypeText.PropString()},
					{Name: "location", DataType: schema.DataTypeText.PropString()},
					{Name: "responsibilities", DataType: schema.DataTypeStringArray.PropString()},
				},
			},
			{
				Name:     "skills",
				DataType: schema.DataTypeStringArray.PropString(),
			},
			{
				Name:     "languages",
				DataType: schema.DataTypeStringArray.PropString(),
			},
			{
				Name:     "projects",
				DataType: schema.DataTypeObjectArray.PropString(),
				NestedProperties: []*models.NestedProperty{
					{Name: "name", DataType: schema.DataTypeText.PropString()},
					{Name: "description", DataType: schema.DataTypeText.PropString()},
					{Name: "technologies", DataType: schema.DataTypeStringArray.PropString()},
					{Name: "link", DataType: schema.DataTypeText.PropString()},
				},
			},
			{
				Name:     "certifications",
				DataType: schema.DataTypeObjectArray.PropString(),
				NestedProperties: []*models.NestedProperty{
					{Name: "name", DataType: schema.DataTypeText.PropString()},
					{Name: "issuer", DataType: schema.DataTypeText.PropString()},
					{Name: "date", DataType: schema.DataTypeText.PropString()},
					{Name: "link", DataType: schema.DataTypeText.PropString()},
				},
			},
			{
				Name:     "publications",
				DataType: schema.DataTypeObjectArray.PropString(),
				NestedProperties: []*models.NestedProperty{
					{Name: "title", DataType: schema.DataTypeText.PropString()},
					{Name: "publisher", DataType: schema.DataTypeText.PropString()},
					{Name: "date", DataType: schema.DataTypeText.PropString()},
					{Name: "link", DataType: schema.DataTypeText.PropString()},
				},
			},
			{
				Name:     "extra",
				DataType: schema.DataTypeText.PropString(),
				ModuleConfig: map[string]any{
					"text2vec-transformers": map[string]any{
						"skip": true, // do NOT vectorize raw JSON
					},
				},
			},
		},
		VectorIndexType: "hnsw",
	}

	if exists, err := c.Schema().
		ClassExistenceChecker().
		WithClassName(className).
		Do(context.Background()); err != nil && !exists {
		if err := c.Schema().
			ClassCreator().
			WithClass(resumeClass).
			Do(context.Background()); err != nil {
			log.Fatalf("error creating class: %v", err.Error())
		}
	}
}
