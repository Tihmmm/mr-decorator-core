package semgrep

import (
	"github.com/Tihmmm/mr-decorator-core/config"
	"github.com/Tihmmm/mr-decorator-core/parser"
	"github.com/Tihmmm/mr-decorator-core/pkg/templater"
)

type SemgrepParser struct {
	cfg *config.SastParserConfig
}

func (p *SemgrepParser) Name() string {
	return "semgrep"
}

func (p *SemgrepParser) Type() string {
	return parser.TypeSast
}

func (p *SemgrepParser) SetConfig(cfg *config.ParserConfig) {
	p.cfg = &cfg.SastParserConfig
}

func (p *SemgrepParser) GetNoteFromReportFile(dir string, subpath string, vulnMgmtId int) (string, error) {
	var genReport parser.GenSast
	return templater.ExecToString(parser.Types[p.Type()], &genReport)
}
