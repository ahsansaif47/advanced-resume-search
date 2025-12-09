package workflows

import (
	"context"
	"time"

	"github.com/ahsansaif47/advanced-resume/internal/parser"
	"github.com/ahsansaif47/advanced-resume/internal/temporalio"
	"github.com/ahsansaif47/advanced-resume/internal/temporalio/activities"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
)

type StoreResumeInputParams struct {
}

type StoreResumeResult struct {
}

func StoreResumeToWeaviate(ctx workflow.Context, data any) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Store Resume Workflow started")

	// Activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
		RetryPolicy:         nil,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var ocrResult string
	if err := workflow.ExecuteActivity(ctx, activities.RunGeminiInference, nil).Get(ctx, &ocrResult); err != nil {
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
	options := client.StartWorkflowOptions{
		ID:        "store-resume-workflow" + uuid.NewString(),
		TaskQueue: temporalio.QueueName,
	}

	r, err := c.ExecuteWorkflow(context.Background(), options, StoreResumeToWeaviate, data)
	if err != nil {
		return "", err
	}
	return r.GetID(), err
}
