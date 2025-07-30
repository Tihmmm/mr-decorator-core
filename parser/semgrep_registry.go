package parser

import (
	"github.com/Tihmmm/mr-decorator-core/config"
	"github.com/Tihmmm/mr-decorator-core/pkg/templater"
)

const formatSemgrep = "semgrep"

type SemgrepParser struct {
	cfg *config.SastParserConfig
}

func (p *SemgrepParser) Name() string {
	return formatSemgrep
}

func (p *SemgrepParser) Type() string {
	return TypeSast
}

func (p *SemgrepParser) SetConfig(cfg *config.ParserConfig) {
	p.cfg = &cfg.SastParserConfig
}

func (p *SemgrepParser) GetNoteFromReportFile(dir string, subpath string, vulnMgmtId int) (string, error) {
	var genReport GenSast
	return templater.ExecToString(Types[p.Type()], &genReport)
}

func init() {
	if isToRegister(formatSemgrep) {
		Register(&SemgrepParser{})
	}
}
