package parser

type cycloneDX struct {
	Vulnerabilities []struct {
		CveId          string `json:"id"`
		Description    string `json:"description"`
		Recommendation string `json:"recommendation"`
		Affects        []struct {
			LibraryName string `json:"ref"`
		} `json:"affects"`
	} `json:"vulnerabilities"`
}

func (dx *cycloneDX) vulnCount() int {
	return len(dx.Vulnerabilities)
}
