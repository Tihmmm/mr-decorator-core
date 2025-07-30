package parser

import (
	"github.com/Tihmmm/mr-decorator-core/config"
)

const formatDependencyTrack = "dependency-track"

type DeptrackParser struct {
	cfg *config.ScaParserConfig
}

func (p *DeptrackParser) Name() string {
	return formatDependencyTrack
}

func (p *DeptrackParser) Type() string {
	return TypeSca
}

func (p *DeptrackParser) SetConfig(cfg *config.ParserConfig) {
	p.cfg = &cfg.ScaParserConfig
}

func (p *DeptrackParser) GetNoteFromReportFile(dir string, subpath string, vulnMgmtId int) (string, error) {
	//TODO implement me
	panic("implement me")
}

func init() {
	if isToRegister(formatDependencyTrack) {
	}
	Register(&DeptrackParser{})
}
