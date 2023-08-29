package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type AgentConfig struct {
	ServerAddress  string `env:"ADDRESS"`
	SignKey        string `env:"KEY"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	RateLimit      int    `env:"RATE_LIMIT"`
	SignKeyByte    []byte
}

const (
	serverAddressDefault  = "localhost:8080"
	reportIntervalDefault = 10
	pollIntervalDefault   = 2
	rateLimitDefault      = 3
)

func NewAgentConfig() *AgentConfig {
	cfg := AgentConfig{}
	cfg.init()

	if cfg.SignKey != "" {
		cfg.SignKeyByte = []byte(cfg.SignKey)
	}

	if err := env.Parse(&cfg); err != nil {
		logrus.Errorf("env parsing error: %v", err)
		return nil
	}
	return &cfg
}

func (c *AgentConfig) init() {
	flag.StringVar(&c.ServerAddress, "a", serverAddressDefault, "Start server address (default - :8080)")
	flag.IntVar(&c.ReportInterval, "r", reportIntervalDefault, "Interval of report metric")
	flag.IntVar(&c.PollInterval, "p", pollIntervalDefault, "Interval of poll metric")
	flag.StringVar(&c.SignKey, "k", "", "Server key")
	flag.IntVar(&c.RateLimit, "l", rateLimitDefault, "Rate limit")
	flag.Parse()
}
