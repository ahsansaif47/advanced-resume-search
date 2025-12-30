package parser

import (
	"encoding/json"
	"log"
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

// func NormalizeKeys(raw map[string]any, aliases map[string]string) map[string]any {
// 	normalized := make(map[string]any)

// 	for key, value := range raw {
// 		if canonical, ok := aliases[key]; ok {
// 			normalized[canonical] = value
// 		} else {
// 			normalized[key] = value
// 		}
// 	}

// 	return normalized
// }

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

// func normalizeSection(
// 	raw any,
// 	aliases map[string]string,
// ) []map[string]any {

// 	var items []map[string]any = []map[string]any{}

// 	switch v := raw.(type) {

// 	case []any:
// 		for _, it := range v {
// 			if m, ok := it.(map[string]any); ok {
// 				items = append(items, normalizeObjectKeys(m, aliases))
// 			}
// 		}

// 	case map[string]any:
// 		items = append(items, normalizeObjectKeys(v, aliases))
// 	}

// 	return items
// }

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

	// normalized := NormalizeKeys(raw, ResumeKeyAliases)
	normalized := NormalizeResumeJSON(raw)

	// log.Printf("Normalized resume data: %+v", normalized)
	prettyData, _ := json.MarshalIndent(normalized, "", " ")
	log.Printf("Normalized resume data (pretty): %s", string(prettyData))

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

	prettyFlat, _ := json.MarshalIndent(flat, "", " ")
	log.Printf("Flattened resume data (pretty): %s", string(prettyFlat))

	var r Resume
	if err := json.Unmarshal(normalizedJson, &r); err != nil {
		return nil, err
	}

	// r.Extra = extra
	return &r, nil
}

// func parseSection[T any](section any) T {
// 	var t T

// 	return t
// }

// Helper function to get string from map with type assertion
// For string, it returns directly
// For []string, it joins them with space and returns
// FIXME: Question: Do I actually need this function?
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

func getStringSlice(m map[string]any, key string) []string {
	if v, ok := m[key]; ok {
		switch v := v.(type) {
		case []string:
			return v
		case []any:
			var result []string
			for _, item := range v {
				if s, ok := item.(string); ok {
					result = append(result, s)
				}
			}
			return result
		case string:
			return []string{v}
		default:
			return []string{}
		}
	}
	return []string{}
}

func prettyPrint(v any) string {
	prettyData, _ := json.MarshalIndent(v, "", " ")
	return string(prettyData)
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

	log.Printf("Normalized education map: %s\n", prettyPrint(normalized["education"]))
	switch normalized["education"].(type) {
	case map[string]any:
		log.Println("Single education entry found (map[string]any).")
		r.Education = append(r.Education, Education{
			Institution: getString(normalized["education"].(map[string]any), "institution"),
			Degree:      getString(normalized["education"].(map[string]any), "degree"),
			Dates:       getString(normalized["education"].(map[string]any), "dates"),
		})

	case []map[string]any:
		log.Println("Array type education entry found ([]map[string]any).")
		for _, e := range normalized["education"].([]map[string]any) {
			r.Education = append(r.Education, Education{
				Institution: getString(e, "institution"),
				Degree:      getString(e, "degree"),
				Dates:       getString(e, "dates"),
			})
		}
	default:
		log.Printf("Education type unsupported: %T", normalized["education"])
	}
	log.Printf("Resume Education: %s", prettyPrint(r.Education))

	switch normalized["work_experience"].(type) {
	case map[string]any:
		log.Println("Single work_experience entry found (map[string]any).")
		work_exp := WorkExp{
			Company:          getString(normalized["work_experience"].(map[string]any), "company"),
			Title:            getString(normalized["work_experience"].(map[string]any), "title"),
			Dates:            getString(normalized["work_experience"].(map[string]any), "dates"),
			Location:         getString(normalized["work_experience"].(map[string]any), "location"),
			Responsibilities: getStringSlice(normalized["work_experience"].(map[string]any), "responsibilities"),
		}
		// FIXME - write a custom empty struct checker
		// if work_exp != (WorkExp{}) {
		// 	r.WorkExperience = append(r.WorkExperience, work_exp)
		// }
		r.WorkExperience = append(r.WorkExperience, work_exp)

	case []map[string]any:
		log.Println("Array type education entry found ([]map[string]any).")
		for _, e := range normalized["work_experience"].([]map[string]any) {
			r.WorkExperience = append(r.WorkExperience, WorkExp{
				Company:          getString(e, "company"),
				Title:            getString(e, "title"),
				Dates:            getString(e, "dates"),
				Location:         getString(e, "location"),
				Responsibilities: getStringSlice(e, "responsibilities"),
			})
		}

	default:
		log.Printf("work_experience type unsupported: %T", normalized["work_experience"])
	}

	switch normalized["projects"].(type) {
	case map[string]any:
		log.Println("Single project entry found (map[string]any).")
		project := normalized["projects"].(map[string]any)
		pr := Project{
			Name:        getString(project, "name"),
			Link:        getString(project, "link"),
			Description: getString(project, "description"),
		}

		if tech, ok := project["technologies"].([]any); ok {
			for _, t := range tech {
				if s, ok := t.(string); ok {
					pr.Technologies = append(pr.Technologies, s)
				}
			}
		}
		r.Projects = append(r.Projects, pr)

	case []map[string]any:
		for _, p := range normalized["projects"].([]map[string]any) {
			pr := Project{
				Name:        getString(p, "name"),
				Link:        getString(p, "link"),
				Description: getString(p, "description"),
			}

			if tech, ok := p["technologies"].([]any); ok {
				for _, t := range tech {
					if s, ok := t.(string); ok {
						pr.Technologies = append(pr.Technologies, s)
					}
				}
			}
			// FIXME - How to check for the empty struct properly?
			// if pr != (Project{}) {
			// 	r.Projects = append(r.Projects, pr)
			// }
			r.Projects = append(r.Projects, pr)
		}

	default:
		log.Printf("Projects type unsupported: %T", normalized["projects"])
	}

	switch normalized["certifications"].(type) {
	case map[string]any:
		log.Println("Single certifications entry found (map[string]any).")

		cert := Certification{
			Name:   getString(normalized["certifications"].(map[string]any), "name"),
			Issuer: getString(normalized["certifications"].(map[string]any), "issuer"),
			Date:   getString(normalized["certifications"].(map[string]any), "date"),
			Link:   getString(normalized["certifications"].(map[string]any), "link"),
		}
		if cert != (Certification{}) {
			r.Certifications = append(r.Certifications, cert)
		}

	case []map[string]any:
		log.Println("Array type certification entry found ([]map[string]any).")
		for _, e := range normalized["certifications"].([]map[string]any) {
			cert := Certification{
				Name:   getString(e, "name"),
				Issuer: getString(e, "issuer"),
				Date:   getString(e, "date"),
				Link:   getString(e, "link"),
			}

			if cert != (Certification{}) {
				r.Certifications = append(r.Certifications, cert)
			}
		}

	default:
		log.Printf("Certifications type unsupported: %T", normalized["certifications"])
	}

	switch normalized["publications"].(type) {
	case map[string]any:
		log.Println("Single publications entry found (map[string]any).")
		pub := Publication{
			Title:     getString(normalized["publications"].(map[string]any), "title"),
			Publisher: getString(normalized["publications"].(map[string]any), "publisher"),
			Date:      getString(normalized["publications"].(map[string]any), "date"),
			Link:      getString(normalized["publications"].(map[string]any), "link"),
		}
		if pub != (Publication{}) {
			r.Publications = append(r.Publications, pub)
		}

	case []map[string]any:
		log.Println("Array type certification entry found ([]map[string]any).")
		for _, e := range normalized["publications"].([]map[string]any) {
			pub := Publication{
				Title:     getString(e, "title"),
				Publisher: getString(e, "publisher"),
				Date:      getString(e, "date"),
				Link:      getString(e, "link"),
			}
			if pub != (Publication{}) {
				r.Publications = append(r.Publications, pub)
			}
		}

	default:
		log.Printf("publications type unsupported: %T", normalized["publications"])
	}

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
