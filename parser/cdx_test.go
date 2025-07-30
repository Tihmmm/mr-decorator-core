package parser

import (
	"github.com/Tihmmm/mr-decorator-core/config"
	"testing"
)

func TestCdxParser_GetNoteFromReportFile(t *testing.T) {
	p := &CdxParser{}
	cfg := &config.ParserConfig{
		ScaParserConfig: config.ScaParserConfig{
			VulnMgmtProjectUrlTmpl: "qwer",
			VulnInstanceTmpl:       "test",
			ReportPath:             "report",
		},
	}
	p.SetConfig(cfg)

	note, err := p.GetNoteFromReportFile("../../test/files/", "trivy.json", 0)
	if err != nil {
		t.Fatalf("GetNoteFromReportFile() error = %v", err)
	}
	t.Logf("note:\n%s\n", note)
}
