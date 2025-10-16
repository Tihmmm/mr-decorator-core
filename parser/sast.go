package parser

import (
	"github.com/Tihmmm/mr-decorator-core/config"
)

const sastVulnsDisplayed = scaVulnsDisplayed

type Sast interface {
	vulnCount() int
	ToGenSast(cfg config.SastParserConfig, vulnMgmtId int) GenSast
}

type GenSast struct {
	VulnMgmtProjectUrl string
	VulnMgmtReportPath string
	HighVulns          []Vulnerability
	CriticalVulns      []Vulnerability
	HighCount          int
	CriticalCount      int
	HcCount            int
}

type Vulnerability struct {
	Name             string
	Location         string
	VulnMgmtInstance string
}

func (sast *GenSast) ApplyLimit() {
	if sast.HighCount > sastVulnsDisplayed {
		sast.HighVulns = sast.HighVulns[:sastVulnsDisplayed]
	}
}
