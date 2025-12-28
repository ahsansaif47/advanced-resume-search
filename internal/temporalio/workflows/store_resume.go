package workflows

import (
	"context"
	"time"

	"github.com/ahsansaif47/advanced-resume/config"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// TODO: Integrate Params and Result into the workflow... best-practice
// TODO: Testing
// TODO: Return the user the result when the workflow completes...

type StoreResumeInputParams struct {
}

type StoreResumeResult struct {
}

func StoreResumeToWeaviate(ctx workflow.Context, data string) (string, error) {
	// logger := workflow.GetLogger(ctx)
	// logger.Info("Store Resume Workflow started")

	// Activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 3,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:    5,
			BackoffCoefficient: 2.0,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var ocrResult string
	if err := workflow.ExecuteActivity(
		ctx,
		"RunGeminiInference",
		data,
	).Get(ctx, &ocrResult); err != nil {
		return "", err
	}

	var inserted_obj_id string
	if err := workflow.ExecuteActivity(ctx,
		"ParseAndStoreData",
		ocrResult,
	).Get(ctx, &inserted_obj_id); err != nil {
		return "", err
	}

	return inserted_obj_id, nil
}

func ExecuteWorkflow_StoreResumeToWeaviate(c client.Client, data string) (string, error) {
	r, err := c.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			TaskQueue: config.ResumeProcessingQueue,
		}, StoreResumeToWeaviate, data)
	if err != nil {
		return "", err
	}
	return r.GetID(), err
}
