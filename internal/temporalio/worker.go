package temporalio

import (
	"log"

	"github.com/ahsansaif47/advanced-resume/internal/temporalio/activities"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var (
	TemporalNamespace = "advanced-resume-parser"
	QueueName         = "resume-processing-queue"
)

func StartWorker() (*client.Client, <-chan error) {
	errCh := make(chan error, 1)

	c, err := client.Dial(client.Options{
		HostPort:  "localhost:7233",
		Namespace: TemporalNamespace,
	})
	if err != nil {
		errCh <- err
		return nil, errCh
	}

	resume_parser := worker.New(c, QueueName, worker.Options{})
	resume_parser.RegisterActivity(activities.RunGeminiInference)
	resume_parser.RegisterActivity(activities.RunOCRDataParsing)
	resume_parser.RegisterActivity(activities.RunStoreResumeDataToWeaviate)

	go func() {
		log.Println("Worker started...")
		if err := resume_parser.Run(worker.InterruptCh()); err != nil {
			errCh <- err
		}
	}()
	return &c, errCh

}
