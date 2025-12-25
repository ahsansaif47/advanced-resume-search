package parser

import (
	"encoding/json"
	"strings"
)

func CleanJSON(input string) string {
	input = strings.TrimSpace(input)

	if strings.HasPrefix(input, "```") {
		idx := strings.Index(input[3:], "\n")

		if idx != -1 {
			input = input[idx+4:]
		}
	}

	input = strings.TrimSuffix(input, "```")

	return input
}

func NormalizeKeys(raw map[string]any, aliases map[string]string) map[string]any {
	normalized := make(map[string]any)

	for key, value := range raw {
		if canonical, ok := aliases[key]; ok {
			normalized[canonical] = value
		} else {
			normalized[key] = value
		}
	}

	return normalized
}

func normalizeObjectArray(v any, allowedKeys map[string]struct{}) []map[string]any {
	var out []map[string]any

	if arr, ok := v.([]any); ok {
		for _, item := range arr {
			if m, ok := item.(map[string]any); ok {
				clean := map[string]any{}

				for k, v := range m {
					if _, ok := allowedKeys[k]; ok {
						clean[k] = v
					}
				}

				out = append(out, clean)
			}
		}
	}
	return out
}

func ParseResume(jsonBytes []byte) (*Resume, error) {
	var raw map[string]any

	if err := json.Unmarshal(jsonBytes, &raw); err != nil {
		return nil, err
	}

	normalized := NormalizeKeys(raw, KeyAliases)

	flat := make(map[string]any)

	if personal_info, ok := normalized["personal_information"].(map[string]any); ok {
		flat["personal_information"] = personal_info
	}

	if experience, ok := normalized["work_experience"].([]any); ok {
		flat["work_experience"] = normalizeObjectArray(
			experience, map[string]struct{}{
				"company":     {},
				"title":       {},
				"dates":       {},
				"description": {},
			})

		// flat["work_experience"] = experience
	}

	if allSkills, ok := normalized["skills"].(map[string]any); ok {
		everySkill := []string{}

		for _, skillSet := range allSkills {
			if skillList, ok := skillSet.([]any); ok {
				for _, skill := range skillList {
					if skillStr, ok := skill.(string); ok {
						everySkill = append(everySkill, skillStr)
					}
				}
			}
		}

		flat["skills"] = everySkill
	}

	if education, ok := normalized["education"]; ok {
		// Getting the normalized data for weaviate
		flat["education"] = normalizeObjectArray(
			education,
			map[string]struct{}{
				"institution": {},
				"degree":      {},
				"dates":       {},
			},
		)

		// flat["education"] = education
	}

	// known := map[string]struct{}{
	// 	"personal_information": {},
	// 	"work_experience":      {},
	// 	"skills":               {},
	// 	"education":            {},
	// 	"languages":            {},
	// }

	// extra := make(map[string]any)
	// for key, val := range normalized {
	// 	if _, ok := known[key]; !ok {
	// 		extra[key] = val
	// 	}
	// }

	normalizedJson, err := json.Marshal(flat)
	if err != nil {
		return nil, err
	}

	var r Resume
	if err := json.Unmarshal(normalizedJson, &r); err != nil {
		return nil, err
	}

	// r.Extra = extra
	return &r, nil
}
