package gemini

import (
	"context"
	"log"
	"os"

	"github.com/ahsansaif47/advanced-resume/config"
	"google.golang.org/genai"
)

type IGeminiClient interface {
	GetResponse(path string) (string, error)
}

type GeminiClient struct {
	client *genai.Client
}

func NewGeminiClient() (IGeminiClient, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: config.GetConfig().GeminiAPIKey,
	})
	if err != nil {
		log.Println("Failed to instantiate genAI client! Err:", err.Error())
		return nil, err
	}

	return &GeminiClient{
		client: client,
	}, nil
}

func (g *GeminiClient) GetResponse(path string) (string, error) {
	ctx := context.Background()

	mimeType := "application/pdf"

	pdfBytes, err := os.ReadFile(path)
	parts := []*genai.Part{
		{
			InlineData: &genai.Blob{
				MIMEType: mimeType,
				Data:     pdfBytes,
			},
		},
		genai.NewPartFromText("Parse the resume to extract the person details like personal information, work experience, skills and all other detials listed on the resume. Return the response structured as json"),
	}

	contents := []*genai.Content{

		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	// Generate content WITH the file
	result, err := g.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite",
		contents,
		nil,
	)
	if err != nil {
		log.Println("Got error from genAI client! Err:", err.Error())
		return "", err
	}

	return result.Text(), nil
}
