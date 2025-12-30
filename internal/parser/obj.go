package parser

type Resume struct {
	PersonalInformation PersonalInfo    `json:"personal_information,omitempty"`
	Summary             string          `json:"summary,omitempty"`
	WorkExperience      []WorkExp       `json:"work_experience,omitempty"`
	Skills              []string        `json:"skills,omitempty"`
	Education           []Education     `json:"education,omitempty"`
	Projects            []Project       `json:"projects,omitempty"`
	Certifications      []Certification `json:"certifications,omitempty"`
	Publications        []Publication   `json:"publications,omitempty"`
	Languages           []string        `json:"languages,omitempty"`
	Extra               string          `json:"resume_extra,omitempty"`
}

type PersonalInfo struct {
	Name     string `json:"name,omitempty"`
	Title    string `json:"title,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Github   string `json:"github,omitempty"`
	LinkedIn string `json:"linkedin,omitempty"`

	Extra map[string]any `json:"info_extra,omitempty"`
}

type WorkExp struct {
	Company          string   `json:"company,omitempty"`
	Location         string   `json:"location,omitempty"`
	Title            string   `json:"title,omitempty"`
	Dates            string   `json:"dates,omitempty"`
	Responsibilities []string `json:"responsibilities,omitempty"`

	// Extra map[string]any `json:"experience_extra,omitempty"`
}

// type Skills struct {
// 	Skills []string `json:"skills,omitempty"`
// }

type Education struct {
	Institution string `json:"institution,omitempty"`
	Location    string `json:"location,omitempty"`
	Degree      string `json:"degree,omitempty"`
	Dates       string `json:"dates,omitempty"`
}

type Project struct {
	Name         string   `json:"name,omitempty"`
	Description  string   `json:"description,omitempty"`
	Technologies []string `json:"technologies,omitempty"`
	Link         string   `json:"link,omitempty"`
	// Extra        map[string]any `json:"project_extra,omitempty"`
}

type Certification struct {
	Name   string `json:"name,omitempty"`
	Issuer string `json:"issuer,omitempty"`
	Date   string `json:"date,omitempty"`
	Link   string `json:"link,omitempty"`
	// Extra  map[string]any `json:"certification_extra,omitempty"`
}

type Publication struct {
	Title     string `json:"title,omitempty"`
	Publisher string `json:"publisher,omitempty"`
	Date      string `json:"date,omitempty"`
	Link      string `json:"link,omitempty"`
	// Extra     map[string]any `json:"publication_extra,omitempty"`
}

var (
	ResumeKeyAliases = map[string]string{
		// Personal Information
		"personal_details": "personal_information",
		"personal_detail":  "personal_information",
		"personalInfo":     "personal_information",
		"personal-info":    "personal_information",
		"info":             "personal_information",

		// Work Experience
		"empolyment_history":      "work_experience",
		"professional_experience": "work_experience",
		"job_history":             "work_experience",
		"workHistory":             "work_experience",
		"experiences":             "work_experience",
		"experience":              "work_experience",
		"jobs":                    "work_experience",
		"work":                    "work_experience",

		// Skills
		"skillset":  "skills",
		"abilities": "skills",
		"tools":     "skills",

		// Education
		"education_history": "education",
		"educationDetails":  "education",
		"academice":         "education",

		// Projects
		"project":               "projects",
		"project_work":          "projects",
		"academic_projects":     "projects",
		"semester_projects":     "projects",
		"final_year_project":    "projects",
		"final_year_projects":   "projects",
		"fyp_personal_projects": "projects",
		"fyp_projects":          "projects",
		"professional_projects": "projects",
		"personal_projects":     "projects",

		// Certifications
		"certifications": "certifications",
		"certs":          "certifications",
		"licenses":       "certifications",

		// Publications
		"publications":          "publications",
		"personal_publications": "publications",
		"papers":                "publications",
		"research":              "publications",
	}

	WorkExpKeyAliases = map[string]string{
		"company":          "company",
		"title":            "title",
		"dates":            "dates",
		"location":         "location",
		"responsibilities": "responsibilities",

		"company_name": "company",
		"employer":     "company",
		"organization": "company",
		"firm":         "company",

		"job_title":   "title",
		"role":        "title",
		"position":    "title",
		"designation": "title",
		"job_role":    "title",

		"duration":    "dates",
		"time_period": "dates",
		"period":      "dates",

		"responsibility":      "responsibilities",
		"role_responsibility": "responsibilities",
		"achievements":        "responsibilities",
		"description":         "responsibilities",
		"contributions":       "responsibilities",
	}

	EducationKeyAliases = map[string]string{

		"institution": "institution",
		"school":      "institution",
		"college":     "institution",
		"university":  "institution",

		"location": "location",

		"degree":        "degree",
		"course":        "degree",
		"program":       "degree",
		"qualification": "degree",

		"dates":           "dates",
		"year":            "dates",
		"duration":        "dates",
		"time_period":     "dates",
		"graduation":      "dates",
		"graduation_year": "dates",
		"period":          "dates",
	}

	ProjectKeyAliases = map[string]string{
		"name":         "name",
		"project_name": "name",
		"title":        "name",

		"details":     "description",
		"info":        "description",
		"description": "description",

		"tech_stack":   "technologies",
		"tools":        "technologies",
		"stack":        "technologies",
		"technologies": "technologies",

		"repository": "link",
		"url":        "link",
		"link":       "link",
		"github":     "link",
	}

	CertificationKeyAliases = map[string]string{
		"name":          "name",
		"cert_name":     "name",
		"certification": "name",
		"title":         "name",

		"issuer":       "issuer",
		"issued_by":    "issuer",
		"issuing_body": "issuer",
		"authority":    "issuer",
		"institution":  "issuer",

		"date":          "date",
		"date_obtained": "date",
		"year":          "date",
		"issued_date":   "date",

		"link": "link",
		"url":  "link",
	}

	PublicationKeyAliases = map[string]string{

		"title":             "title",
		"publication_title": "title",

		"publisher_name": "publisher",
		"publisher":      "publisher",

		"date_published": "date",
		"year":           "date",
		"published_year": "date",
		"date":           "date",

		"url":           "link",
		"link":          "link",
		"article_link":  "link",
		"articles_link": "link",
		"article_url":   "link",
		"articles_url":  "link",
	}

	PersonalInfoKeyAliases = map[string]string{

		"name":      "name",
		"full_name": "name",

		"title":    "title",
		"role":     "title",
		"position": "title",

		"email":         "email",
		"mail":          "email",
		"email_address": "email",

		"phone":          "phone",
		"phone_number":   "phone",
		"contact":        "phone",
		"contact_number": "phone",
		"contact_no":     "phone",

		"github":      "github",
		"github_url":  "github",
		"github_link": "github",

		"linkedin":      "linkedin",
		"linkedin_url":  "linkedin",
		"linkedin_link": "linkedin",

		"address":          "address",
		"current_location": "address",
		"current_address":  "address",
	}
)

var SectionAliases = map[string]map[string]string{
	"personal_information": PersonalInfoKeyAliases,
	"work_experience":      WorkExpKeyAliases,
	"education":            EducationKeyAliases,
	"projects":             ProjectKeyAliases,
	"certifications":       CertificationKeyAliases,
	"publications":         PublicationKeyAliases,
}
