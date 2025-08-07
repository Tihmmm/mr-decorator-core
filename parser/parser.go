package parser

import (
	"github.com/Tihmmm/mr-decorator-core/config"
	"slices"
)

type Parser interface {
	Name() string
	Type() string
	SetConfig(cfg *config.ParserConfig)
	GetNoteFromReportFile(dir string, subpath string, vulnMgmtId int) (string, error)
}

const (
	formatAny = "*"
)

func isToRegister(format string) bool {
	return config.RegisteredParsers == nil || slices.Contains(config.RegisteredParsers, formatAny) || config.RegisteredParsers[0] == "" || slices.Contains(config.RegisteredParsers, format)
}
