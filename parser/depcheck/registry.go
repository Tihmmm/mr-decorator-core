package depcheck

import (
	"fmt"
	"github.com/Tihmmm/mr-decorator-core/config"
	"github.com/Tihmmm/mr-decorator-core/parser"
	"github.com/Tihmmm/mr-decorator-core/pkg/file"
	"github.com/Tihmmm/mr-decorator-core/pkg/templater"
	"log"
)

type DepCheckParser struct {
	cfg *config.ScaParserConfig
}

func (p *DepCheckParser) Name() string {
	return "dependency-check"
}

func (p *DepCheckParser) Type() string {
	return parser.TypeSca
}

func (p *DepCheckParser) SetConfig(cfg *config.ParserConfig) {
	p.cfg = &cfg.ScaParserConfig
}

func (p *DepCheckParser) GetNoteFromReportFile(dir string, subpath string, vulnMgmtId int) (string, error) {
	var depcheck dependencyCheck
	if err := file.ParseJsonFile(dir, subpath, &depcheck); err != nil {
		log.Printf("error parsing dependency check file: %v\n", err)
		return "", err
	}

	var genReport parser.GenSca
	parseGenReport(vulnMgmtId, p.cfg, &depcheck, &genReport)

	genReport.ApplyLimit()

	return templater.ExecToString(parser.Types[p.Type()], &genReport)
}

func parseGenReport(vulnMgmtId int, cfg *config.ScaParserConfig, dc *dependencyCheck, dest *parser.GenSca) {
	dest.Count = dc.vulnCount()
	baseUrl := fmt.Sprintf(cfg.VulnMgmtProjectUrlTmpl, vulnMgmtId)
	dest.VulnMgmtProjectUrl = baseUrl
	for _, v := range dc.Dependencies {
		for _, vuln := range v.Vulnerabilities {
			cve := parser.Cve{
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
	parser.Register(&DepCheckParser{})
}
