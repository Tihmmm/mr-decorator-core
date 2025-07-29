package fpr

import (
	"fmt"
	"github.com/Tihmmm/mr-decorator-core/config"
	"github.com/Tihmmm/mr-decorator-core/parser"
	"github.com/Tihmmm/mr-decorator-core/pkg/templater"
	"log"
)

type FprParser struct {
	cfg *config.SastParserConfig
}

func (p *FprParser) Name() string {
	return "fpr"
}

func (p *FprParser) Type() string {
	return parser.TypeSast
}

func (p *FprParser) SetConfig(cfg *config.ParserConfig) {
	p.cfg = &cfg.SastParserConfig
}

func (p *FprParser) GetNoteFromReportFile(dir string, _ string, vulnMgmtId int) (string, error) {
	var fprr fpr
	if err := ParseFprFile(dir, &fprr); err != nil {
		log.Printf("error parsing fpr file: %v\n", err)
		return "", err
	}

	var genReport parser.GenSast
	parseGenReport(vulnMgmtId, p.cfg, &fprr, &genReport)

	genReport.ApplyLimit()

	return templater.ExecToString(parser.Types[p.Type()], genReport)
}

func parseGenReport(vulnMgmtId int, cfg *config.SastParserConfig, fprr *fpr, dest *parser.GenSast) {
	dest.HcCount = fprr.vulnCount()
	dest.HighCount = fprr.highCount
	dest.CriticalCount = fprr.criticalCount
	baseUrl := fmt.Sprintf(cfg.VulnMgmtProjectUrlTmpl, vulnMgmtId)
	dest.VulnMgmtProjectUrl = baseUrl
	for _, v := range fprr.highRecords {
		highVulns := parser.Vulnerability{
			Name:             v.category,
			Location:         v.path,
			VulnMgmtInstance: baseUrl + fmt.Sprintf(cfg.VulnInstanceTmpl, v.sscVulnInstance),
		}
		dest.HighVulns = append(dest.HighVulns, highVulns)
	}
	for _, v := range fprr.criticalRecords {
		criticalVulns := parser.Vulnerability{
			Name:             v.category,
			Location:         v.path,
			VulnMgmtInstance: baseUrl + fmt.Sprintf(cfg.VulnInstanceTmpl, v.sscVulnInstance),
		}
		dest.CriticalVulns = append(dest.CriticalVulns, criticalVulns)
	}
	dest.VulnMgmtReportPath = baseUrl + cfg.ReportPath
}

func init() {
	parser.Register(&FprParser{})
}
