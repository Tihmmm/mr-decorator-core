package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Parser       ParserConfig     `yaml:"parser"`
	GitlabClient BaseClientConfig `yaml:"client"`
	Decorator    DecoratorConfig  `yaml:"decorator"`
}

type DecoratorConfig struct {
	ArtifactDownloadMaxRetries int `yaml:"artifact_download_max_retries" default:"3"`
	ArtifactDownloadRetryDelay int `yaml:"artifact_download_retry_delay" default:"2"` // seconds
}

type BaseClientConfig struct {
	Ip             string `yaml:"ip"`
	Host           string `yaml:"host"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
}

type ParserConfig struct {
	ScaParserConfig   ScaParserConfig  `yaml:"sca"`
	SastParserConfig  SastParserConfig `yaml:"sast"`
	registeredParsers []string         `yaml:"registered_parsers"`
}

type ScaParserConfig struct {
	VulnMgmtProjectUrlTmpl string `yaml:"vuln_mgmt_project_url_tmpl"`
	VulnInstanceTmpl       string `yaml:"vuln_instance_tmpl"`
	ReportPath             string `yaml:"report_path"`
}

type SastParserConfig struct {
	VulnMgmtProjectUrlTmpl string `yaml:"vuln_mgmt_project_url_tmpl"` // e.g. https://fortify-ssc.company.com/html/ssc/version/%d
	VulnInstanceTmpl       string `yaml:"vuln_instance_tmpl"`         // e.g. audit?q=instance_id%3A
	ReportPath             string `yaml:"report_path"`                // e.g. audit?q=analysis_type%3Asca
}

func NewConfig(path string) (Config, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, errors.New(fmt.Sprintf("Error reading config.yml: %s\n", err))
	}

	var cfg Config
	buf := bytes.NewBuffer(configBytes)
	dec := yaml.NewDecoder(buf)
	if err := dec.Decode(&cfg); err != nil {
		return Config{}, errors.New(fmt.Sprintf("Error parsing config.yml: %s\n", err))
	}

	return cfg, nil
}
