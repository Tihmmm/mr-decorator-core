package parser

import (
	"github.com/Tihmmm/mr-decorator-core/config"
)

const scaVulnsDisplayed = 10

type Sca interface {
	vulnCount() int
	ToGenSca(cfg config.ScaParserConfig, vulnMgmtId int) GenSca
}

type GenSca struct {
	VulnMgmtProjectUrl string
	VulnMgmtReportPath string
	Cves               []Cve
	Count              int
}

type Cve struct {
	Id               string
	LibraryName      string
	Description      string
	Recommendations  string
	VulnMgmtInstance string
}

func (sca *GenSca) ApplyLimit() {
	if sca.Count > scaVulnsDisplayed {
		sca.Cves = sca.Cves[:scaVulnsDisplayed]
	}
}
