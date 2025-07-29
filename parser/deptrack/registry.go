package deptrack

import (
	"github.com/Tihmmm/mr-decorator-core/config"
	"github.com/Tihmmm/mr-decorator-core/parser"
)

type DeptrackParser struct {
	cfg *config.ScaParserConfig
}

func (p *DeptrackParser) Name() string {
	return "dependency-track"
}

func (p *DeptrackParser) Type() string {
	return parser.TypeSca
}

func (p *DeptrackParser) SetConfig(cfg *config.ParserConfig) {
	p.cfg = &cfg.ScaParserConfig
}

func (p *DeptrackParser) GetNoteFromReportFile(dir string, subpath string, vulnMgmtId int) (string, error) {
	//TODO implement me
	panic("implement me")
}
