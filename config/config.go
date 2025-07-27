package config

import (
	"bytes"
	"github.com/goccy/go-yaml"
	"log"
	"os"
)

type Config struct {
	Parser       ParserConfig       `yaml:"parser"`
	Server       ServerConfig       `yaml:"server"`
	GitlabClient GitlabClientConfig `yaml:"gitlab_client"`
}

type ParserConfig struct {
	ScaParserConfig  ScaParserConfig  `yaml:"sca"`
	SastParserConfig SastParserConfig `yaml:"sast"`
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
