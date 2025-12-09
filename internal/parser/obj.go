package parser

type Resume struct {
	PersonalInformation PersonalInfo `json:"personal_information,omitempty"`
	Summary             string       `json:"summary,omitempty"`
	WorkExperience      []WorkExp    `json:"work_experience,omitempty"`
	Skills              []string     `json:"skills,omitempty"`
	Education           []Education  `json:"education,omitempty"`

	Extra map[string]any `json:"-,omitempty"`
}

type PersonalInfo struct {
	Name     string `json:"name,omitempty"`
	Title    string `json:"title,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Github   string `json:"github,omitempty"`
	LinkedIn string `json:"linkedin,omitempty"`

	Extra map[string]any `json:"-,omitempty"`
}

type WorkExp struct {
	Company          string   `json:"company,omitempty"`
	Location         string   `json:"location,omitempty"`
	Title            string   `json:"title,omitempty"`
	Dates            string   `json:"dates,omitempty"`
	Responsibilities []string `json:"responsibilities,omitempty"`

	Extra map[string]any `json:"-,omitempty"`
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

var KeyAliases = map[string]string{
	// Personal Information
	"personal_details": "personal_information",
	"personal_detail":  "personal_information",
	"personalInfo":     "personal_information",
	"personal-info":    "personal_information",
	"info":             "personal_information",

	// Work Experience
	"empolyment_history": "work_experience",
	"job_history":        "work_experience",
	"workHistory":        "work_experience",
	"experiences":        "work_experience",
	"experience":         "work_experience",
	"jobs":               "work_experience",
	"work":               "work_experience",

	// Skills
	"skillset":  "skills",
	"abilities": "skills",

	// Education
	"education_history": "education",
	"educationDetails":  "education",
	"academice":         "education",
}
