package config

import (
	"bytes"
	"github.com/goccy/go-yaml"
	"log"
	"os"
)

type Config struct {
	Server       ServerConfig       `yaml:"server"`
	GitlabClient GitlabClientConfig `yaml:"gitlab_client"`
}

type ServerConfig struct {
	Port string `yaml:"port" default:"3000"`
}
type GitlabClientConfig struct {
	Ip   string `yaml:"ip"`
	Host string `yaml:"host"`
}

func NewGeneralConfig(path string) Config {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config.yml: %s\n", err)
	}

	var cfg Config
	buf := bytes.NewBuffer(configBytes)
	dec := yaml.NewDecoder(buf)
	if err := dec.Decode(&cfg); err != nil {
		log.Fatalf("Error parsing config.yml: %s\n", err)
	}

	return cfg
}
