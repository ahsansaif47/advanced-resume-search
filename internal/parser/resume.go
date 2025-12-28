package parser

import (
	"encoding/json"
	"fmt"
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
				"company":          {},
				"title":            {},
				"dates":            {},
				"description":      {},
				"location":         {},
				"responsibilities": {},
				"achievements":     {},
			})

		// flat["work_experience"] = experience
	}

	everySkill := []string{}
	switch normalized["skills"].(type) {
	case map[string]any:
		for _, skillSet := range normalized["skills"].(map[string]any) {
			if skillList, ok := skillSet.([]any); ok {
				for _, skill := range skillList {
					if skillStr, ok := skill.(string); ok {
						everySkill = append(everySkill, skillStr)
					}
				}
			}
		}

		flat["skills"] = everySkill

	case []any:
		for _, skill := range normalized["skills"].([]any) {
			if skillStr, ok := skill.(string); ok {
				everySkill = append(everySkill, skillStr)
			}
		}

		flat["skills"] = everySkill
	}

	if languages, ok := normalized["languages"].([]any); ok {
		out := []string{}
		for _, lang := range languages {
			if lang_str, ok := lang.(string); ok {
				out = append(out, lang_str)
			}
		}
		flat["languages"] = out
	}

	if education, ok := normalized["education"]; ok {
		// Getting the normalized data for weaviate
		flat["education"] = normalizeObjectArray(
			education,
			map[string]struct{}{
				"institution": {},
				"degree":      {},
				"dates":       {},
				"location":    {},
			},
		)
	}

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

func getString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		switch v := v.(type) {
		case []string:
			var desc strings.Builder
			for _, str := range v {
				desc.WriteString(str + " ")
			}
			return desc.String()
		case string:
			return v
		default:
			return ""
		}

	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func ParseResumeUpdated(jsonBytes []byte) (*Resume, error) {
	var raw map[string]any
	if err := json.Unmarshal(jsonBytes, &raw); err != nil {
		return nil, err
	}

	// Normalize top-level keys
	normalized := NormalizeKeys(raw, KeyAliases)

	r := &Resume{
		Extra: make(map[string]any),
	}

	if pinfo, ok := normalized["personal_information"].(map[string]any); ok {
		r.PersonalInformation = PersonalInfo{
			Name:     getString(pinfo, "name"),
			Email:    getString(pinfo, "email"),
			Phone:    getString(pinfo, "phone"),
			Title:    getString(pinfo, "title"),
			LinkedIn: getString(pinfo, "linkedin"),
			Github:   getString(pinfo, "github"),
		}
	}

	if summary, ok := normalized["summary"].(string); ok {
		r.Summary = summary
	}

	switch normalized["skills"].(type) {
	case map[string]any:
		for _, skillSet := range normalized["skills"].(map[string]any) {
			if skillList, ok := skillSet.([]any); ok {
				for _, skill := range skillList {
					if skillStr, ok := skill.(string); ok {
						r.Skills = append(r.Skills, skillStr)
					}
				}
			}
		}

	case []any:
		for _, skill := range normalized["skills"].([]any) {
			if skillStr, ok := skill.(string); ok {
				r.Skills = append(r.Skills, skillStr)
			}
		}
	default:
		// Do nothing
	}

	if edu, ok := normalized["education"].([]any); ok {
		for _, e := range edu {
			if m, ok := e.(map[string]any); ok {
				r.Education = append(r.Education, Education{
					Institution: getString(m, "institution"),
					Degree:      getString(m, "degree"),
					Dates:       getString(m, "dates"),
				})
			}
		}
	}

	if exp, ok := normalized["work_experience"].([]any); ok {
		for _, e := range exp {
			if m, ok := e.(map[string]any); ok {
				r.WorkExperience = append(r.WorkExperience, WorkExp{
					Company:  firstNonEmpty(getString(m, "company"), getString(m, "organization")),
					Title:    firstNonEmpty(getString(m, "title"), getString(m, "designation"), getString(m, "role")),
					Dates:    firstNonEmpty(getString(m, "dates"), getString(m, "tenure")),
					Location: firstNonEmpty(getString(m, "location"), getString(m, "address")),
					// Responsibilities: firstNonEmpty(getString(m, "responsibilities"), getString(m, "tasks")),
				})
			}
		}
	}

	if projs, ok := normalized["projects"].([]any); ok {
		for _, p := range projs {
			if m, ok := p.(map[string]any); ok {
				pr := Project{
					Name:        getString(m, "name"),
					Link:        getString(m, "link"),
					Description: getString(m, "description"),
				}

				if tech, ok := m["technologies"].([]any); ok {
					for _, t := range tech {
						if s, ok := t.(string); ok {
							pr.Technologies = append(pr.Technologies, s)
						}
					}
				}

				r.Projects = append(r.Projects, pr)
			}
		}
	}

	if certs, ok := normalized["certifications"].([]any); ok {
		for _, c := range certs {
			if m, ok := c.(map[string]any); ok {
				cert := Certification{
					Name: firstNonEmpty(
						getString(m, "name"), getString(m, "course"), getString(m, "title"),
					),
					Issuer: firstNonEmpty(
						getString(m, "issuer"), getString(m, "institution"), getString(m, "organization"),
					),
					Date: firstNonEmpty(
						getString(m, "date"), fmt.Sprint(m["year"]),
					),
				}

				// IMPORTANT: skip empty objects
				if cert.Name != "" || cert.Issuer != "" || cert.Date != "" {
					r.Certifications = append(r.Certifications, cert)
				}
			}
		}
	}

	if pubs, ok := normalized["publications"].([]any); ok {
		for _, p := range pubs {
			if m, ok := p.(map[string]any); ok {
				r.Publications = append(r.Publications, Publication{
					Title:     getString(m, "title"),
					Publisher: getString(m, "publisher"),
					Date:      getString(m, "date"),
					Link:      getString(m, "link"),
				})
			}
		}
	}

	// ---- Everything else → Extra ----
	known := map[string]struct{}{
		"summary":              {},
		"personal_information": {},
		"skills":               {},
		"education":            {},
		"work_experience":      {},
		"projects":             {},
		"languages":            {},
		"certifications":       {},
		"publications":         {},
	}

	for k, v := range normalized {
		if _, ok := known[k]; !ok {
			r.Extra[k] = v
		}
	}

	return r, nil
}
