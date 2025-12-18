package controllers

import (
	"fmt"
	"mime/multipart"

	"github.com/ahsansaif47/advanced-resume/integrations/gemini"
	"github.com/ahsansaif47/advanced-resume/internal/storage/weaviate"
	"github.com/ahsansaif47/advanced-resume/internal/temporalio/workflows"
	"github.com/ahsansaif47/advanced-resume/utils"
	"github.com/gofiber/fiber/v2"
	"go.temporal.io/sdk/client"
)

type IWeaviateService interface {
	AddResumeToDB(*fiber.Ctx, *multipart.FileHeader) (int, string, error)
	BatchUploadResume(batchResume []map[string]any) (int, error)
	VectorSearch(query string) (int, any, error)
}

type WeaviateService struct {
	repo           weaviate.IWeaviateRepository
	GeminiClient   gemini.IGeminiClient
	TemporalClient client.Client
}

// FIXME: Use Interfaces for temporal client!!
func NewWeaviateService(
	repo weaviate.IWeaviateRepository,
	genAIclient gemini.IGeminiClient,
	temporalClient client.Client,
) IWeaviateService {
	return &WeaviateService{
		repo:           repo,
		GeminiClient:   genAIclient,
		TemporalClient: temporalClient,
	}
}

// TODO: Move these to constants package..
// NOTE: Move the fmt.Errorf statements to errors package
// NOTE: Create both packages
var (
	InternalServerError = fiber.StatusInternalServerError
	StatusAccepted      = fiber.StatusAccepted
	StatusOK            = fiber.StatusOK
)

func (s *WeaviateService) AddResumeToDB(ctx *fiber.Ctx, resumeFile *multipart.FileHeader) (int, string, error) {

	// Sanitize FileName
	sanitizedFName := utils.SanitizeFileName(resumeFile.Filename)

	// File Upload to resources
	err := ctx.SaveFile(resumeFile, fmt.Sprintf("../../resources/pdfs/%s", sanitizedFName))
	if err != nil {
		return InternalServerError, "", fmt.Errorf("Error saving file: %s", err.Error())
	}

	resumePath := fmt.Sprintf("../../resources/pdfs/%s", sanitizedFName)

	// // Temp Workflow

	// // Activity: Extract Resume Info
	// resumeData, err := s.GeminiClient.GetResponse(resumePath) // 6s
	// if err != nil {
	// 	return InternalServerError, "", fmt.Errorf("Error in OCR: %s", err.Error())
	// }

	// // Clean data
	// cleanedData := parser.CleanJSON(resumeData)

	// // Parse into obj
	// data, err := parser.ParseResume([]byte(cleanedData))
	// if err != nil {
	// 	return InternalServerError, "", fmt.Errorf("Error parsing resume: %s", err.Error())
	// }

	// var bytesData []byte
	// if bytesData, err = json.MarshalIndent(data, "", ""); err != nil {
	// 	return InternalServerError, "", fmt.Errorf("Error marshalling data: %s", err.Error())
	// }

	// // splits := strings.Split(resumePath, "/")
	// // fPath := fmt.Sprintf("../../tmp/%s", strings.Replace(fmt.Sprintf("%s", splits[len(splits)-1]), ".pdf", ".json", -1))

	// // // Create or open file
	// // file, err := os.Create(fPath)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }
	// // defer file.Close()

	// // // Write bytes directly
	// // _, err = file.Write(bytesData)
	// // if err != nil {
	// // 	log.Fatal(err)
	// // }

	// var requiredData map[string]any
	// if err := json.Unmarshal(bytesData, &requiredData); err != nil {
	// 	return InternalServerError, "", fmt.Errorf("Error unmarshalling data: %s", err.Error())
	// }

	// id, err := s.repo.AddResumeToDB("resume", requiredData)
	// if err != nil {
	// 	return InternalServerError, "", fmt.Errorf("Error uploading resume: %s", err.Error())
	// }

	workflow_id, err := workflows.ExecuteWorkflow_StoreResumeToWeaviate(s.TemporalClient, resumePath)

	return StatusAccepted, workflow_id, nil
}

func (s *WeaviateService) BatchUploadResume(batchResume []map[string]any) (int, error) {
	err := s.repo.BatchUploadResume("resume", batchResume)
	if err != nil {
		return InternalServerError, fmt.Errorf("Error while batch uploading: %s", err.Error())
	}

	return StatusAccepted, nil
}

func (s *WeaviateService) VectorSearch(query string) (int, any, error) {
	data, err := s.repo.VectorSearch("resume", query)
	if err != nil {
		return InternalServerError, nil, fmt.Errorf("Error searchinh: %s", err.Error())
	}
	return StatusOK, data, nil
}
