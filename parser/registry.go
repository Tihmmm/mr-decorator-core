package parser

import (
	"github.com/Tihmmm/mr-decorator-core/errors"
	"log"
	"sync"
)

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

	return nil, &errors.FormatError{}
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
