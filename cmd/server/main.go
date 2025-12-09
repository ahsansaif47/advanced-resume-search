package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ahsansaif47/advanced-resume/integrations/gemini"
	"github.com/ahsansaif47/advanced-resume/internal/parser"
	"github.com/ahsansaif47/advanced-resume/internal/storage/weaviate"
)

var className = "Resume"

func main() {
	c, err := gemini.NewGeminiClient()
	if err != nil {
		log.Panic("Err creating genAI client. Err: ", err.Error())
	}
	log.Println("Client created")
	weaviateClient := weaviate.ConnectWeaviate()
	weaviate.CreateSchema(weaviateClient, className)

	repo := weaviate.NewWeviateRepository(context.Background(), weaviateClient)

	resume_files := []string{"/home/ahsansaif/projects/advanced-resume/resources/pdfs/AhsanResume202507.pdf"}
	for _, resume := range resume_files {
		res, err := c.GetResponse(resume)
		if err != nil {
			log.Fatal("Error: ", err.Error())
		}
		cleanedData := parser.CleanJSON(res)

		data, err := parser.ParseResume([]byte(cleanedData))
		if err != nil {
			log.Fatal("Error: ", err.Error())
		}

		var bytesData []byte

		if bytesData, err = json.MarshalIndent(data, "", " "); err != nil {
			log.Fatal("Error: ", err.Error())
		}

		splits := strings.Split(resume, "/")
		fPath := fmt.Sprintf("../../tmp/%s", strings.Replace(fmt.Sprintf("%s", splits[len(splits)-1]), ".pdf", ".json", -1))

		// Create or open file
		file, err := os.Create(fPath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Write bytes directly
		_, err = file.Write(bytesData)
		if err != nil {
			log.Fatal(err)
		}

		var requiredData map[string]any
		if err := json.Unmarshal(bytesData, &requiredData); err != nil {
			log.Fatal("Failed to UnMarshal json response")
		}

		id, err := repo.AddResumeToDB(className, requiredData)
		if err != nil {
			fmt.Printf("%s", err.Error())
			log.Println("Resume not added")
		}

		fmt.Println("Id for the inserted resume is: ", id)

	}

	result, err := repo.VectorSearch(className, "Get me resumes for Golang developers")
	if err != nil {
		return
	}
	rawResult := result["Get"].(map[string]any)["Resume"].([]interface{})

	resumes := []parser.Resume{}
	for _, raw := range rawResult {
		b, _ := json.Marshal(raw)

		var resume parser.Resume

		if err := json.Unmarshal(b, &resume); err != nil {
			log.Println("Failed to parse into resume object")
		}
		bytz, _ := json.MarshalIndent(resume, "", " ")
		fmt.Printf("%s\n\n", bytz)
		resumes = append(resumes, resume)
	}
	log.Println("Printing data objectss. Len: ", len(resumes))
}
