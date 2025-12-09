package workflows

import (
	"time"

	"github.com/ahsansaif47/advanced-resume/internal/parser"
	"github.com/ahsansaif47/advanced-resume/internal/temporalio/activities"
	"go.temporal.io/sdk/workflow"
)

func StoreResumeToWeaviate(ctx workflow.Context, data any) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Store Resume Workflow started")

	// Activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
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
