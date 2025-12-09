package workflows

import (
	"encoding/json"
	"fmt"

	"github.com/ahsansaif47/advanced-resume/internal/parser"
)

func StoreResumeToWeaviate() {
	// Activity: Extract Resume Info
	resumeData, err := s.GeminiClient.GetResponse(resumePath) // 6s
	if err != nil {
		return InternalServerError, "", fmt.Errorf("Error in OCR: %s", err.Error())
	}

	// Clean data
	cleanedData := parser.CleanJSON(resumeData)

	// Parse into obj
	data, err := parser.ParseResume([]byte(cleanedData))
	if err != nil {
		return InternalServerError, "", fmt.Errorf("Error parsing resume: %s", err.Error())
	}

	var bytesData []byte
	if bytesData, err = json.MarshalIndent(data, "", ""); err != nil {
		return InternalServerError, "", fmt.Errorf("Error marshalling data: %s", err.Error())
	}

	// splits := strings.Split(resumePath, "/")
	// fPath := fmt.Sprintf("../../tmp/%s", strings.Replace(fmt.Sprintf("%s", splits[len(splits)-1]), ".pdf", ".json", -1))

	// // Create or open file
	// file, err := os.Create(fPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// // Write bytes directly
	// _, err = file.Write(bytesData)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	var requiredData map[string]any
	if err := json.Unmarshal(bytesData, &requiredData); err != nil {
		return InternalServerError, "", fmt.Errorf("Error unmarshalling data: %s", err.Error())
	}

	id, err := s.repo.AddResumeToDB("resume", requiredData)
	if err != nil {
		return InternalServerError, "", fmt.Errorf("Error uploading resume: %s", err.Error())
	}

}
