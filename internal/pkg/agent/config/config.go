package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type AgentConfig struct {
	ServerAddress  string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"` //тесты не проходят с duration
	PollInterval   int    `env:"POLL_INTERVAL"`   //тесты не проходят с duration
}

func NewAgentConfig() *AgentConfig {
	cfg := AgentConfig{}
	cfg.init()

	if err := env.Parse(&cfg); err != nil {
		logrus.Errorf("env parsing error: %v", err)
		return nil
	}
	return &cfg
}

func (c *AgentConfig) init() {
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "Start server address (default - :8080)")
	flag.IntVar(&c.ReportInterval, "r", 10, "Interval of report metric")
	flag.IntVar(&c.PollInterval, "p", 2, "Interval of poll metric")
	flag.Parse()
}
