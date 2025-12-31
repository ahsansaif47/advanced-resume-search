package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ahsansaif47/advanced-resume/config"
	"github.com/ahsansaif47/advanced-resume/internal/parser"
	"github.com/ahsansaif47/advanced-resume/internal/storage/weaviate"
)

func (a *Activities) ParseAndStoreData(ctx context.Context, resumeText string) (string, error) {
	repo := weaviate.NewWeviateRepository(ctx, a.WeaviateClient)

	// Clean data
	cleanedData := parser.CleanJSON(resumeText)

	// log.Println("Cleaned resume data:", cleanedData)

	// Parse into obj
	data, err := parser.ParseResumeUpdated([]byte(cleanedData))
	// prettyStr, err := json.MarshalIndent(data, "", " ")
	// if err != nil {
	// 	return "", fmt.Errorf("Error marshalling data: %s", err.Error())``
	// }
	// log.Println("Parsed resume data:", data)

	if err != nil {
		return "", fmt.Errorf("Error parsing resume: %s", err.Error())
	}

	var bytesData []byte
	if bytesData, err = json.MarshalIndent(data, "", " "); err != nil {
		return "", fmt.Errorf("Error marshalling data: %s", err.Error())
	}

	log.Println("Parsed resume data:", string(bytesData))

	var resumeMapData map[string]any
	if err := json.Unmarshal(bytesData, &resumeMapData); err != nil {
		return "", fmt.Errorf("Error unmarshalling data: %s", err.Error())
	}

	// NOTE: Sanitize the map before inserting data into weaviate..
	// resumeMapData = sanitizeMap(resumeMapData)

	id, err := repo.AddResumeToDB(config.ClassName, resumeMapData)
	if err != nil {
		return "", fmt.Errorf("Error uploading resume: %s", err.Error())
	}
	// utils.SaveResumeDataJson(id, bytesData)

	return id, nil
}
