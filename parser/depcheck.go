package parser

type dependencyCheck struct {
	Dependencies []struct {
		LibraryName     string `json:"fileName"`
		Vulnerabilities []struct {
			CveId       string `json:"name"`
			Description string `json:"description"`
		} `json:"vulnerabilities,omitempty"`
	} `json:"dependencies"`
}

func (dc *dependencyCheck) vulnCount() int {
	var count int
	for i := range dc.Dependencies {
		count += len(dc.Dependencies[i].Vulnerabilities)
	}

	return count
}
