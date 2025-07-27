package parser

import (
	"fmt"
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

	//p.SetConfig()

	registry[key] = p
}

func Get(format string) (Parser, error) {
	mu.RLock()
	defer mu.RUnlock()

	if p, ok := registry[format]; ok {
		return p, nil
	}

	return nil, fmt.Errorf("no parser registered for format %q\n", format)
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
