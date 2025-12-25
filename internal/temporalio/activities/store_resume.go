package activities

import (
	"regexp"
	"strings"
)

func sanitizeKey(key string) string {
	if key == "" {
		return ""
	}

	validKeyRegex := regexp.MustCompile(`[^_0-9A-Za-z]`)

	key = strings.TrimSpace(key)
	key = validKeyRegex.ReplaceAllString(key, "_")

	if key[0] >= '0' && key[0] <= '9' {
		key = "_" + key
	}

	return key
}

func sanitizeValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		return sanitizeMap(val)
	case []any:
		out := make([]any, 0, len(val))
		for _, item := range val {
			out = append(out, sanitizeValue(item))
		}
		return out
	default:
		return val
	}

}

func sanitizeMap(m map[string]any) map[string]any {
	out := make(map[string]any)
	for k, v := range m {
		out[sanitizeKey(k)] = sanitizeValue(v)
	}
	return out
}

// func (a *Activities) RunStoreResumeDataToWeaviate(ctx context.Context, resume parser.Resume) (id string, err error) {

// 	repo := weaviate.NewWeviateRepository(context.Background(), weaviate.ConnectWeaviate())

// 	var bytesData []byte
// 	if bytesData, err = json.MarshalIndent(resume, "", ""); err != nil {
// 		return "", fmt.Errorf("Error marshalling data: %s", err.Error())
// 	}

// 	var resumeMapData map[string]any
// 	if err := json.Unmarshal(bytesData, &resumeMapData); err != nil {
// 		return "", fmt.Errorf("Error unmarshalling data: %s", err.Error())
// 	}

// 	// NOTE: Sanitize the map before inserting data into weaviate..
// 	resumeMapData = sanitizeMap(resumeMapData)

// 	id, err = repo.AddResumeToDB("resume", resumeMapData)
// 	if err != nil {
// 		return "", fmt.Errorf("Error uploading resume: %s", err.Error())
// 	}
// 	utils.SaveResumeDataJson(id, bytesData)

// 	return id, nil
// }
