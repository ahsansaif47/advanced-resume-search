package activities

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/genai"
)

// var genClient gemini.IGeminiClient

func (a *Activities) RunGeminiInference(ctx context.Context, path string) (string, error) {

	// var err error
	resumeData, err := a.GenAIClient.GetResponse(path) // 6s
	if err != nil {
		if errors.Is(err, genai.APIError{}) {
			// FIXME: Handle rate limit exceeded error
			// If err is rate limit exceeded, stop further retrying the request

			// "message": "Error runnig ocr! Err: Error 429, Message: You exceeded your current quota, please check your plan and billing details. For more information on this error, head to: https://ai.google.dev/gemini-api/docs/rate-limits. To monitor your current usage, head to: https://ai.dev/usage?tab=rate-limit.
			// * Quota exceeded for metric: generativelanguage.googleapis.com/generate_content_free_tier_requests, limit: 20, model: gemini-2.5-flash-lite
			// Please retry in 25.618799928s., Status: RESOURCE_EXHAUSTED, Details: [map[@type:type.googleapis.com/google.rpc.Help links:[map[description:Learn more about Gemini API quotas url:https://ai.google.dev/gemini-api/docs/rate-limits]]] map[@type:type.googleapis.com/google.rpc.QuotaFailure violations:[map[quotaDimensions:map[location:global model:gemini-2.5-flash-lite] quotaId:GenerateRequestsPerDayPerProjectPerModel-FreeTier quotaMetric:generativelanguage.googleapis.com/generate_content_free_tier_requests quotaValue:20]]] map[@type:type.googleapis.com/google.rpc.RetryInfo retryDelay:25s]]",

			return "", nil
		}
		return "", fmt.Errorf("Error runnig ocr! Err: %s", err.Error())
	}

	return resumeData, nil
}
