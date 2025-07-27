package parser

import (
	"errors"
	"github.com/Tihmmm/mr-decorator-core/config"
	"log"
	"sync"
)

const (
	scaVulnsDisplayed  = 10
	sastVulnsDisplayed = scaVulnsDisplayed
)

type Parser interface {
	Name() string
	Type() string
	SetConfig(cfg *config.ParserConfig)
	GetNoteFromReportFile(dir string, subpath string, vulnMgmtId int) (string, error)
}

var (
	mu       sync.RWMutex
	registry = make(map[string]Parser)
)

func Register(p Parser) {
	mu.Lock()
	defer mu.Unlock()

	key := p.Name()
	if _, exists := registry[key]; exists {
		log.Fatalf("Parser already registered: %s", p.Name())
	}

	registry[key] = p
}

func Get(format string) (Parser, error) {
	mu.RLock()
	defer mu.RUnlock()

	if p, ok := registry[format]; ok {
		return p, nil
	}

	return nil, errors.New("parser not registered")
}

func List() []string {
	mu.RLock()
	defer mu.RUnlock()

	keys := make([]string, 0, len(registry))
	for k := range registry {
		keys = append(keys, k)
	}

	return keys
}
