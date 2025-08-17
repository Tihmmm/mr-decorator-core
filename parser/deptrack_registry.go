package parser

import (
	"fmt"

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
	return "", fmt.Errorf("not implemented: %s", formatDependencyTrack)
}

func init() {
	Register(&DeptrackParser{})
}
