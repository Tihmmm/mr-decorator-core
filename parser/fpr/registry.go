package fpr

import (
	"fmt"
	"github.com/Tihmmm/mr-decorator-core/config"
	"github.com/Tihmmm/mr-decorator-core/models"
	"github.com/Tihmmm/mr-decorator-core/parser"
	"github.com/Tihmmm/mr-decorator-core/pkg/templater"
	"log"
)

type FprParser struct {
	cfg *config.SastParserConfig
}

func (p FprParser) Name() string {
	return "fpr"
}

func (p FprParser) Type() string {
	return models.FormatSast
}

func (p FprParser) SetConfig(cfg *config.ParserConfig) {
	p.cfg = &cfg.SastParserConfig
}

func (p FprParser) GetNoteFromReportFile(dir string, _ string, vulnMgmtId int) (string, error) {
	var fprr fpr
	if err := ParseFprFile(dir, &fprr); err != nil {
		log.Printf("error parsing fpr file: %v\n", err)
		return "", err
	}

	var genReport parser.GenSast
	parseGenReport(vulnMgmtId, p.cfg, &fprr, &genReport)

	genReport.ApplyLimit()

	return templater.ExecToString(parser.BaseTemplateSast, &genReport)
}

func parseGenReport(vulnMgmtId int, cfg *config.SastParserConfig, fprr *fpr, dest *parser.GenSast) {
	var genReport parser.GenSast

	genReport.HcCount = fprr.vulnCount()
	genReport.HighCount = fprr.highCount
	genReport.CriticalCount = fprr.criticalCount
	baseUrl := fmt.Sprintf(cfg.VulnMgmtProjectUrlTmpl, vulnMgmtId)
	genReport.VulnMgmtProjectUrl = baseUrl
	for _, v := range fprr.highRecords {
		highVulns := parser.Vulnerability{
			Name:             v.category,
			Location:         v.path,
			VulnMgmtInstance: baseUrl + fmt.Sprintf(cfg.VulnInstanceTmpl, v.sscVulnInstance),
		}
		genReport.HighVulns = append(genReport.HighVulns, highVulns)
	}
	for _, v := range fprr.criticalRecords {
		criticalVulns := parser.Vulnerability{
			Name:             v.category,
			Location:         v.path,
			VulnMgmtInstance: baseUrl + fmt.Sprintf(cfg.VulnInstanceTmpl, v.sscVulnInstance),
		}
		genReport.CriticalVulns = append(genReport.CriticalVulns, criticalVulns)
	}
	genReport.VulnMgmtReportPath = baseUrl + cfg.ReportPath

	dest = &genReport
}

func init() {
	parser.Register(FprParser{})
}
