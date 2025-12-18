package workflows

import (
	"context"
	"time"

	"github.com/ahsansaif47/advanced-resume/config"
	"github.com/ahsansaif47/advanced-resume/internal/parser"
	"github.com/ahsansaif47/advanced-resume/internal/temporalio/activities"
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
	if err := workflow.ExecuteActivity(ctx, activities.RunGeminiInference, data).Get(ctx, &ocrResult); err != nil {
		return "", err
	}

	var resume parser.Resume
	if err := workflow.ExecuteActivity(ctx, activities.RunOCRDataParsing, ocrResult).Get(ctx, &resume); err != nil {
		return "", err
	}

	var inserted_obj_id string
	if err := workflow.ExecuteActivity(ctx, activities.RunStoreResumeDataToWeaviate, resume).Get(ctx, &inserted_obj_id); err != nil {
		return "", err
	}

	return inserted_obj_id, nil
}

func ExecuteWorkflow_StoreResumeToWeaviate(c client.Client, data string) (string, error) {
	// options := client.StartWorkflowOptions{
	// 	// ID:                    "store-resume-workflow" + uuid.NewString(),
	// 	TaskQueue: config.QueueName,
	// 	// WorkflowIDReusePolicy: enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
	// }

	r, err := c.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			TaskQueue: config.QueueName,
		}, StoreResumeToWeaviate, data)
	if err != nil {
		return "", err
	}
	return r.GetID(), err
}
