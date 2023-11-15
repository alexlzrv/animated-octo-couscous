package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type AgentConfig struct {
	ServerAddress  string `env:"ADDRESS" json:"server_address"`
	SignKey        string `env:"KEY"`
	ReportInterval int    `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   int    `env:"POLL_INTERVAL" json:"poll_interval"`
	RateLimit      int    `env:"RATE_LIMIT"`
	PublicKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigPath     string `env:"CONFIG"`
	SignKeyByte    []byte
}

const (
	serverAddressDefault  = "localhost:8080"
	reportIntervalDefault = 10
	pollIntervalDefault   = 2
	rateLimitDefault      = 3
)

func NewAgentConfig() (*AgentConfig, error) {
	cfg := AgentConfig{}
	cfg.init()

	if cfg.SignKey != "" {
		cfg.SignKeyByte = []byte(cfg.SignKey)
	}

	if cfg.ConfigPath != "" {
		cfgJSON, err := readConfigFile(cfg.ConfigPath)
		if err != nil {
			return cfgJSON, err
		}
	}

	if err := env.Parse(&cfg); err != nil {
		logrus.Errorf("env parsing error: %v", err)
		return nil, err
	}
	return &cfg, nil
}

func (c *AgentConfig) init() {
	flag.StringVar(&c.ServerAddress, "a", serverAddressDefault, "Start server address (default - :8080)")
	flag.IntVar(&c.ReportInterval, "r", reportIntervalDefault, "Interval of report metric")
	flag.IntVar(&c.PollInterval, "p", pollIntervalDefault, "Interval of poll metric")
	flag.StringVar(&c.SignKey, "k", "", "Server key")
	flag.IntVar(&c.RateLimit, "l", rateLimitDefault, "Rate limit")
	flag.StringVar(&c.PublicKey, "-crypto-key", "", "Public key path")
	flag.StringVar(&c.ConfigPath, "c", "", "Path to config file")
	flag.StringVar(&c.ConfigPath, "config", "", "Path to config file (the same as -c)")
	flag.Parse()
}

func readConfigFile(path string) (cfg *AgentConfig, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
