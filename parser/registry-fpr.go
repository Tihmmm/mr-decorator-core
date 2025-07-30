package parser

import (
	"fmt"
	"github.com/Tihmmm/mr-decorator-core/config"
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
	return TypeSast
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

	var genReport GenSast
	parseFprGenReport(vulnMgmtId, p.cfg, &fprr, &genReport)

	genReport.ApplyLimit()

	return templater.ExecToString(Types[p.Type()], &genReport)
}

func parseFprGenReport(vulnMgmtId int, cfg *config.SastParserConfig, fprr *fpr, dest *GenSast) {
	dest.HcCount = fprr.vulnCount()
	dest.HighCount = fprr.highCount
	dest.CriticalCount = fprr.criticalCount
	baseUrl := fmt.Sprintf(cfg.VulnMgmtProjectUrlTmpl, vulnMgmtId)
	dest.VulnMgmtProjectUrl = baseUrl
	for _, v := range fprr.highRecords {
		highVulns := Vulnerability{
			Name:             v.category,
			Location:         v.path,
			VulnMgmtInstance: baseUrl + fmt.Sprintf(cfg.VulnInstanceTmpl, v.sscVulnInstance),
		}
		dest.HighVulns = append(dest.HighVulns, highVulns)
	}
	for _, v := range fprr.criticalRecords {
		criticalVulns := Vulnerability{
			Name:             v.category,
			Location:         v.path,
			VulnMgmtInstance: baseUrl + fmt.Sprintf(cfg.VulnInstanceTmpl, v.sscVulnInstance),
		}
		dest.CriticalVulns = append(dest.CriticalVulns, criticalVulns)
	}
	dest.VulnMgmtReportPath = baseUrl + cfg.ReportPath
}

func Init(cfg *config.SastParserConfig) {
	Register(
		&FprParser{
			cfg: cfg,
		},
	)
}
