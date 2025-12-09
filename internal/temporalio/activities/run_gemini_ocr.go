package activities

import (
	"context"
	"fmt"

	"github.com/ahsansaif47/advanced-resume/integrations/gemini"
)

var genClient gemini.IGeminiClient

func RunGeminiInference(ctx context.Context, path string) (string, error) {

	var err error
	genClient, err = gemini.NewGeminiClient()
	if err != nil {
		return "", err
	}
	resumeData, err := genClient.GetResponse(path) // 6s
	if err != nil {
		return "", fmt.Errorf("Error runnig ocr! Err: %s", err.Error())
	}

	return resumeData, nil
}
