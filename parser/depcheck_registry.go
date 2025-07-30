package parser

import (
	"errors"
	"fmt"
	"github.com/Tihmmm/mr-decorator-core/config"
	"github.com/Tihmmm/mr-decorator-core/pkg/file"
	"github.com/Tihmmm/mr-decorator-core/pkg/templater"
)

const formatDependencyCheck = "dependency-check"

type DepCheckParser struct {
	cfg *config.ScaParserConfig
}

func (p *DepCheckParser) Name() string {
	return formatDependencyCheck
}

func (p *DepCheckParser) Type() string {
	return TypeSca
}

func (p *DepCheckParser) SetConfig(cfg *config.ParserConfig) {
	p.cfg = &cfg.ScaParserConfig
}

func (p *DepCheckParser) GetNoteFromReportFile(dir string, subpath string, vulnMgmtId int) (string, error) {
	var depcheck dependencyCheck
	if err := file.ParseJsonFile(dir, subpath, &depcheck); err != nil {
		return "", errors.New(fmt.Sprintf("error parsing dependency check file: %v\n", err))
	}

	var genReport GenSca
	parseDepcheckGenReport(vulnMgmtId, p.cfg, &depcheck, &genReport)

	genReport.ApplyLimit()

	return templater.ExecToString(Types[p.Type()], &genReport)
}

func parseDepcheckGenReport(vulnMgmtId int, cfg *config.ScaParserConfig, dc *dependencyCheck, dest *GenSca) {
	dest.Count = dc.vulnCount()
	baseUrl := fmt.Sprintf(cfg.VulnMgmtProjectUrlTmpl, vulnMgmtId)
	dest.VulnMgmtProjectUrl = baseUrl
	for _, v := range dc.Dependencies {
		for _, vuln := range v.Vulnerabilities {
			cve := Cve{
				Id:          vuln.CveId,
				LibraryName: v.LibraryName,
				Description: vuln.Description,
			}
			dest.Cves = append(dest.Cves, cve)
		}
	}
	dest.VulnMgmtReportPath = baseUrl + cfg.ReportPath
}

func init() {
	if isToRegister(formatDependencyCheck) {
		Register(&DepCheckParser{})
	}
}
