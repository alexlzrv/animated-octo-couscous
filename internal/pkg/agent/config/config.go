package config

import (
	"flag"
	"os"
	"strconv"
)

type AgentConfig struct {
	ServerAddress  string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"` //тесты не проходят с duration
	PollInterval   int    `env:"POLL_INTERVAL"`   //тесты не проходят с duration
}

func NewAgentConfig() *AgentConfig {
	return &AgentConfig{}
}

func Init(c *AgentConfig) {
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "Start server address (default - :8080)")
	flag.IntVar(&c.ReportInterval, "r", 10, "Interval of report metric")
	flag.IntVar(&c.PollInterval, "p", 2, "Interval of poll metric")
	flag.Parse()

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		c.ServerAddress = envServerAddress
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		c.ReportInterval, _ = strconv.Atoi(envReportInterval)
	}
	if envPollInterval := os.Getenv("POLL_INTERVAL"); envPollInterval != "" {
		c.PollInterval, _ = strconv.Atoi(envPollInterval)
	}
}
