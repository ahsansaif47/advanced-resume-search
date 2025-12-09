package activities

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ahsansaif47/advanced-resume/internal/parser"
	"github.com/ahsansaif47/advanced-resume/internal/storage/weaviate"
	"github.com/ahsansaif47/advanced-resume/utils"
)

// FIXME: Do we need to do it like this or should we make a global instance at repository.go ??
var repo weaviate.IWeaviateRepository

func RunStoreResumeDataToWeaviate(ctx context.Context, resume parser.Resume) (id string, err error) {

	repo = weaviate.NewWeviateRepository(context.Background(), weaviate.ConnectWeaviate())

	var bytesData []byte
	if bytesData, err = json.MarshalIndent(resume, "", ""); err != nil {
		return "", fmt.Errorf("Error marshalling data: %s", err.Error())
	}

	var resumeMapData map[string]any
	if err := json.Unmarshal(bytesData, &resumeMapData); err != nil {
		return "", fmt.Errorf("Error unmarshalling data: %s", err.Error())
	}

	id, err = repo.AddResumeToDB("resume", resumeMapData)
	if err != nil {
		return "", fmt.Errorf("Error uploading resume: %s", err.Error())
	}
	utils.SaveResumeDataJson(id, bytesData)

	return id, nil
}
