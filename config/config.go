package config

import (
	"bytes"
	"github.com/goccy/go-yaml"
	"log"
	"os"
)

type Config struct {
	Decorator    DecoratorConfig    `yaml:"decorator"`
	Server       ServerConfig       `yaml:"server"`
	GitlabClient GitlabClientConfig `yaml:"gitlab_client"`
	Parser       ParserConfig       `yaml:"parser"`
}

type DecoratorConfig struct {
	ArtifactDownloadMaxRetries int `yaml:"artifact_download_max_retries" default:"3"`
	ArtifactDownloadRetryDelay int `yaml:"artifact_download_retry_delay" default:"2"` // seconds
}

type ServerConfig struct {
	Port             string   `yaml:"port" default:"3000"`
	RateLimit        int      `yaml:"rate_limit" default:"3"`
	ApiKey           string   `yaml:"api_key"`
	SupportedFormats []string `yaml:"supported_formats"`
}
type GitlabClientConfig struct {
	Ip   string `yaml:"ip"`
	Host string `yaml:"host"`
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

func NewConfig(path string) (Config, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading config.yml: %s\n", err)
		return Config{}, err
	}

	var cfg Config
	buf := bytes.NewBuffer(configBytes)
	dec := yaml.NewDecoder(buf)
	if err := dec.Decode(&cfg); err != nil {
		log.Printf("Error parsing config.yml: %s\n", err)
		return Config{}, err
	}

	return cfg, nil
}
