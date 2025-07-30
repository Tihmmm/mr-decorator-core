package parser

import (
	"fmt"
	"github.com/Tihmmm/mr-decorator-core/config"
	"github.com/Tihmmm/mr-decorator-core/pkg/file"
	"github.com/Tihmmm/mr-decorator-core/pkg/templater"
	"log"
)

// CdxParser for json CycloneDX 1.6 SBOMs as outputted by trivy
type CdxParser struct {
	cfg *config.ScaParserConfig
}

func (p *CdxParser) Name() string {
	return "cyclonedx"
}

func (p *CdxParser) Type() string {
	return TypeSca
}

func (p *CdxParser) SetConfig(cfg *config.ParserConfig) {
	p.cfg = &cfg.ScaParserConfig
}

func (p *CdxParser) GetNoteFromReportFile(dir string, subpath string, vulnMgmtId int) (string, error) {
	var cdx cycloneDX
	if err := file.ParseJsonFile(dir, subpath, &cdx); err != nil {
		log.Printf("error parsing cyclonedx file: %v\n", err)
		return "", err
	}

	var genReport GenSca
	parseCdxGenReport(vulnMgmtId, p.cfg, &cdx, &genReport)

	genReport.ApplyLimit()

	return templater.ExecToString(Types[p.Type()], &genReport)
}

func parseCdxGenReport(vulnMgmtId int, cfg *config.ScaParserConfig, cdx *cycloneDX, dest *GenSca) {
	dest.Count = cdx.vulnCount()
	if cfg.VulnMgmtProjectUrlTmpl != "" {
		dest.VulnMgmtProjectUrl = fmt.Sprintf(cfg.VulnMgmtProjectUrlTmpl, vulnMgmtId)
	}
	for _, vuln := range cdx.Vulnerabilities {
		cve := Cve{
			Id:              vuln.CveId,
			LibraryName:     vuln.Affects[0].LibraryName,
			Description:     vuln.Description,
			Recommendations: vuln.Recommendation,
		}
		dest.Cves = append(dest.Cves, cve)
	}
	dest.VulnMgmtReportPath = dest.VulnMgmtProjectUrl + cfg.ReportPath
}

func init() {
	Register(&CdxParser{})
}
