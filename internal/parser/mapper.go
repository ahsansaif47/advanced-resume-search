package parser

import "strings"

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

// Helper function to get string from map with type assertion
// For string, it returns directly
// FIXME -For []string, it joins them with space and returns
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

func educationMapper(edu map[string]any) Education {
	return Education{
		Institution: getString(edu, "institution"),
		Degree:      getString(edu, "degree"),
		Dates:       getString(edu, "dates"),
	}
}

func experienceMapper(workXp map[string]any) WorkExp {
	return WorkExp{
		Company:          getString(workXp, "company"),
		Title:            getString(workXp, "title"),
		Dates:            getString(workXp, "dates"),
		Location:         getString(workXp, "location"),
		Responsibilities: getStringSlice(workXp, "responsibilities"),
	}
}

func projectsMapper(pr map[string]any) Project {
	project := Project{
		Name:        getString(pr, "name"),
		Link:        getString(pr, "link"),
		Description: getString(pr, "description"),
	}

	if tech, ok := pr["technologies"].([]any); ok {
		for _, t := range tech {
			if s, ok := t.(string); ok {
				project.Technologies = append(project.Technologies, s)
			}
		}
	}
	return project
}
func certificationsMapper(cert map[string]any) Certification {
	return Certification{
		Name:   getString(cert, "name"),
		Issuer: getString(cert, "issuer"),
		Date:   getString(cert, "date"),
		Link:   getString(cert, "link"),
	}
}

func publicationsMapper(pub map[string]any) Publication {
	return Publication{
		Title:     getString(pub, "title"),
		Publisher: getString(pub, "publisher"),
		Date:      getString(pub, "date"),
		Link:      getString(pub, "link"),
	}
}
