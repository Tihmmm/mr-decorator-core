package parser

import (
	"log"
	"testing"

	"github.com/Tihmmm/mr-decorator-core/config"
)

func TestDepCheckParser_GetNoteFromReportFile(t *testing.T) {
	p := &DepCheckParser{}
	cfg := &config.ParserConfig{
		ScaParserConfig: config.ScaParserConfig{},
	}
	p.SetConfig(cfg)

	note, err := p.GetNoteFromReportFile("../test/files/", "dependency-check-report.json", 0)
	if err != nil {
		t.Fatalf("GetNoteFromReportFile() error = %v", err)
	}

	log.Printf("note:\n%s", note)
}
