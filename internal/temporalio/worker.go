package temporalio

import (
	"context"
	"log"

	"github.com/ahsansaif47/advanced-resume/config"
	"github.com/ahsansaif47/advanced-resume/internal/temporalio/activities"
	"github.com/ahsansaif47/advanced-resume/internal/temporalio/workflows"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func StartWorker() (client.Client, <-chan error) {
	errCh := make(chan error, 1)

	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
		// Namespace: config.TemporalNamespace,
	})
	if err != nil {
		errCh <- err
		return nil, errCh
	}

	acts, err := activities.NewActivities(context.Background())
	if err != nil {
		errCh <- err
	}

	resume_parser := worker.New(c, config.QueueName, worker.Options{})

	resume_parser.RegisterActivityWithOptions(
		acts.RunGeminiInference,
		activity.RegisterOptions{
			Name: "RunGeminiInference",
		},
	)

	resume_parser.RegisterActivityWithOptions(
		acts.ParseAndStoreData,
		activity.RegisterOptions{
			Name: "ParseAndStoreData",
		},
	)
	// resume_parser.RegisterActivity(activities.RunStoreResumeDataToWeaviate)

	resume_parser.RegisterWorkflow(workflows.StoreResumeToWeaviate)

	go func() {
		log.Println("Worker started...")
		if err := resume_parser.Run(worker.InterruptCh()); err != nil {
			errCh <- err
		}
	}()
	return c, errCh

}
