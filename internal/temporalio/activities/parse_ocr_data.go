package activities

import (
	"context"
	"fmt"

	"github.com/ahsansaif47/advanced-resume/internal/parser"
)

func RunOCRDataParsing(ctx context.Context, resumeText string) (*parser.Resume, error) {
	// Clean data
	cleanedData := parser.CleanJSON(resumeText)

	// Parse into obj
	data, err := parser.ParseResume([]byte(cleanedData))
	if err != nil {
		return nil, fmt.Errorf("Error parsing resume: %s", err.Error())
	}

	return data, nil
}
