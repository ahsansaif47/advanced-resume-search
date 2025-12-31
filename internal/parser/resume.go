package parser

import (
	"encoding/json"
	"log"
	"strings"
)

// FIXME - Fix duplicate keys and normalization logic...

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

func normalizeAny(v any, aliases map[string]string) any {
	switch val := v.(type) {

	case map[string]any:
		out := make(map[string]any)
		for k, v2 := range val {
			key := k
			if aliases != nil {
				if mapped, ok := aliases[k]; ok {
					key = mapped
				}
			}
			out[key] = normalizeAny(v2, aliases)
		}
		return out

	case []any:
		for i := range val {
			val[i] = normalizeAny(val[i], aliases)
		}
		return val

	default:
		return v
	}
}

func normalizeMap(
	m map[string]any,
	aliases map[string]string,
) map[string]any {

	out := make(map[string]any)
	extra := make(map[string]any)

	for k, v := range m {
		if canon, ok := aliases[k]; ok {
			out[canon] = normalizeAny(v, nil)
		} else {
			extra[k] = normalizeAny(v, nil)
		}
	}

	// Attach extras safely
	if len(extra) > 0 {
		out["_extra"] = extra
	}

	return out
}

func normalizeSection(v any, aliases map[string]string) any {
	switch val := v.(type) {

	case []any:
		out := []map[string]any{}
		for _, item := range val {
			if m, ok := item.(map[string]any); ok {
				norm := normalizeMap(m, aliases)
				if len(norm) > 0 {
					out = append(out, norm)
				}
			}
		}
		return out

	case map[string]any:
		return normalizeMap(val, aliases)

	default:
		return v
	}
}

func NormalizeResumeJSON(raw map[string]any) map[string]any {
	out := make(map[string]any)

	for k, v := range raw {
		// Normalize top-level key
		canonKey := k
		if mapped, ok := ResumeKeyAliases[k]; ok {
			canonKey = mapped
		}

		// Apply section-specific normalization
		if aliasMap, ok := SectionAliases[canonKey]; ok && aliasMap != nil {
			out[canonKey] = normalizeSection(v, aliasMap)
		} else {
			// Generic recursion
			out[canonKey] = normalizeAny(v, nil)
		}
	}

	return out
}

func prettyPrint(v any) string {
	prettyData, _ := json.MarshalIndent(v, "", " ")
	return string(prettyData)
}

func parseSection[T any](
	raw any,
	sectionName string,
	mapper func(map[string]any) T,
	appendFn func(T),
) {
	switch val := raw.(type) {
	case map[string]any:
		log.Printf("Single %s entry found", sectionName)
		item := mapper(val)
		appendFn(item)

	case []map[string]any:
		log.Printf("Array type %s entry found", sectionName)
		for _, e := range val {
			item := mapper(e)
			appendFn(item)
		}
	default:
		log.Printf("%s type unsupported: %T", sectionName, raw)
	}
}

func ParseResumeUpdated(jsonBytes []byte) (*Resume, error) {
	var raw map[string]any
	if err := json.Unmarshal(jsonBytes, &raw); err != nil {
		return nil, err
	}

	normalized := NormalizeResumeJSON(raw)
	prettyData, _ := json.MarshalIndent(normalized, "", " ")
	log.Printf("Normalized resume data (pretty): %s", string(prettyData))

	r := &Resume{}

	if pinfo, ok := normalized["personal_information"].(map[string]any); ok {
		log.Printf("Personal Information: %s\n", prettyPrint(pinfo))

		info := PersonalInfo{
			Name:     getString(pinfo, "name"),
			Email:    getString(pinfo, "email"),
			Phone:    getString(pinfo, "phone"),
			Title:    getString(pinfo, "title"),
			LinkedIn: getString(pinfo, "linkedin"),
			Github:   getString(pinfo, "github"),
		}

		r.PersonalInformation = info
		log.Printf("Personal Information: %s", prettyPrint(r.PersonalInformation))
	} else {
		log.Println("No personal information found.")
	}

	if summary, ok := normalized["summary"].(string); ok {
		log.Printf("Summery: %s\n", summary)
		r.Summary = summary
	} else {
		log.Println("No summary found.")
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
		log.Printf("Skills type unsupported: %T", normalized["skills"])
	}

	parseSection(normalized["education"], "education", educationMapper, func(e Education) {
		if e != (Education{}) {
			r.Education = append(r.Education, e)
		}
	})

	parseSection(normalized["work_experience"], "work_experience", experienceMapper, func(ex WorkExp) {
		// FIXME - Custom validator for []string in struct can not compare directly
		r.WorkExperience = append(r.WorkExperience, ex)
	})

	parseSection(normalized["projects"], "projects", projectsMapper, func(p Project) {
		// FIXME - Custom validator for []string in struct can not compare directly
		r.Projects = append(r.Projects, p)
	})

	parseSection(normalized["certifications"], "certifications", certificationsMapper, func(c Certification) {
		if c != (Certification{}) {
			r.Certifications = append(r.Certifications, c)
		}
		// r.Certifications = append(r.Certifications, c)
	})

	parseSection(normalized["publications"], "publications", publicationsMapper, func(p Publication) {
		if p != (Publication{}) {
			r.Publications = append(r.Publications, p)
		}
		// r.Publications = append(r.Publications, p)
	})

	// ---- Everything else to Extra ----
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

	extraMap := make(map[string]any)
	for k, v := range normalized {
		if _, ok := known[k]; !ok {
			extraMap[k] = v
		}
	}

	extraBytes, err := json.Marshal(extraMap)
	if err != nil {
		return r, err
	}
	r.Extra = string(extraBytes)

	prettyResume, _ := json.MarshalIndent(r, "", " ")
	log.Printf("Pretty resume: %s", string(prettyResume))
	return r, nil
}
