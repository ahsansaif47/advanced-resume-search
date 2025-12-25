package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ahsansaif47/advanced-resume/internal/parser"
	"github.com/ahsansaif47/advanced-resume/internal/storage/weaviate"
)

func (a *Activities) ParseAndStoreData(ctx context.Context, resumeText string) (string, error) {
	// Clean data
	cleanedData := parser.CleanJSON(resumeText)

	// Parse into obj
	data, err := parser.ParseResume([]byte(cleanedData))
	if err != nil {
		return "", fmt.Errorf("Error parsing resume: %s", err.Error())
	}

	repo := weaviate.NewWeviateRepository(ctx, a.WeaviateClient)

	var bytesData []byte
	if bytesData, err = json.MarshalIndent(data, "", ""); err != nil {
		return "", fmt.Errorf("Error marshalling data: %s", err.Error())
	}

	var resumeMapData map[string]any
	if err := json.Unmarshal(bytesData, &resumeMapData); err != nil {
		return "", fmt.Errorf("Error unmarshalling data: %s", err.Error())
	}

	// NOTE: Sanitize the map before inserting data into weaviate..
	resumeMapData = sanitizeMap(resumeMapData)

	id, err := repo.AddResumeToDB("resume", resumeMapData)
	if err != nil {
		return "", fmt.Errorf("Error uploading resume: %s", err.Error())
	}
	// utils.SaveResumeDataJson(id, bytesData)

	return id, nil
}
