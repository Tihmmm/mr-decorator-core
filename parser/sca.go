package parser

import (
	"github.com/Tihmmm/mr-decorator-core/config"
)

type Sca interface {
	vulnCount() int
	ToGenSca(cfg config.ScaParserConfig, vulnMgmtId int) GenSca
}

type GenSca struct {
	Count              int
	Cves               []Cve
	VulnMgmtProjectUrl string
	VulnMgmtReportPath string
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
