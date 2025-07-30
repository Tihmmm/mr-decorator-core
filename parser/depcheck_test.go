package parser

import (
	"github.com/Tihmmm/mr-decorator-core/config"
	"log"
	"testing"
)

func TestDepCheckParser_GetNoteFromReportFile(t *testing.T) {
	p := &DepCheckParser{}
	cfg := &config.ParserConfig{
		ScaParserConfig: config.ScaParserConfig{
			VulnMgmtProjectUrlTmpl: "qwer",
			VulnInstanceTmpl:       "test",
			ReportPath:             "report",
		},
	}
	p.SetConfig(cfg)

	note, err := p.GetNoteFromReportFile("../../test/files/", "depcheck.json", 0)
	if err != nil {
		t.Fatalf("GetNoteFromReportFile() error = %v", err)
	}

	log.Printf("note:\n%s", note)
}
