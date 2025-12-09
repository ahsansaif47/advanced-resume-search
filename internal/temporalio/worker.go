package temporalio

import (
	"log"

	"github.com/ahsansaif47/advanced-resume/internal/temporalio/activities"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func start_worker() *client.Client {
	c, err := client.Dial(client.Options{
		HostPort:  "",
		Namespace: "Advanced-resume-parser",

		// Credentials: client.NewMTLSCredentials(tls.Certificate{}),
	})
	if err != nil {
		log.Fatalln("Unable to create client:", err)
	}

	process_resume_worker := worker.New(c, "resume-processing-queue", worker.Options{})

	process_resume_worker.RegisterActivity(activities.RunGeminiInference)
	process_resume_worker.RegisterActivity(activities.RunOCRDataParsing)
	process_resume_worker.RegisterActivity(activities.RunStoreResumeDataToWeaviate)

	log.Println("Worker started...")
	err = process_resume_worker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker:", err)
	}

	return &c
}
