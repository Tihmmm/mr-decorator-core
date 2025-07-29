package cdx

import (
	"fmt"
	"github.com/Tihmmm/mr-decorator-core/config"
	"github.com/Tihmmm/mr-decorator-core/parser"
	"github.com/Tihmmm/mr-decorator-core/pkg/file"
	"github.com/Tihmmm/mr-decorator-core/pkg/templater"
	"log"
)

type CdxParser struct {
	cfg *config.ScaParserConfig
}

func (p CdxParser) Name() string {
	return "cyclonedx"
}

func (p CdxParser) Type() string {
	return parser.TypeSca
}

func (p CdxParser) SetConfig(cfg *config.ParserConfig) {
	p.cfg = &cfg.ScaParserConfig
}

func (p CdxParser) GetNoteFromReportFile(dir string, subpath string, vulnMgmtId int) (string, error) {
	var cdx cycloneDX
	if err := file.ParseJsonFile(dir, subpath, &cdx); err != nil {
		log.Printf("error parsing cyclonedx file: %v\n", err)
		return "", err
	}

	var genReport parser.GenSca
	parseGenReport(vulnMgmtId, p.cfg, &cdx, &genReport)

	genReport.ApplyLimit()

	return templater.ExecToString(parser.Types[p.Type()], genReport)
}

func parseGenReport(vulnMgmtId int, cfg *config.ScaParserConfig, cdx *cycloneDX, dest *parser.GenSca) {
	dest.Count = cdx.vulnCount()
	baseUrl := fmt.Sprintf(cfg.VulnMgmtProjectUrlTmpl, vulnMgmtId)
	dest.VulnMgmtProjectUrl = baseUrl
	for _, vuln := range cdx.Vulnerabilities {
		cve := parser.Cve{
			Id:              vuln.CveId,
			LibraryName:     vuln.Affects[0].LibraryName,
			Description:     vuln.Description,
			Recommendations: vuln.Recommendation,
		}
		dest.Cves = append(dest.Cves, cve)
	}
	dest.VulnMgmtReportPath = baseUrl + cfg.ReportPath
}

func init() {
	parser.Register(CdxParser{})
}
