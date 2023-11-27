package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
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
	PublicKeyPath  string `env:"CRYPTO_KEY" json:"crypto_key"`
	PublicKey      *rsa.PublicKey
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

	if err := env.Parse(&cfg); err != nil {
		logrus.Errorf("env parsing error: %v", err)
		return nil, err
	}

	if cfg.ConfigPath != "" {
		cfgJSON, err := readConfigFile(cfg.ConfigPath)
		if err != nil {
			return cfgJSON, err
		}
	}

	if cfg.PublicKeyPath != "" {
		publicKey, err := cfg.getPublicKey()
		if err != nil {
			logrus.Errorf("error with get public key: %v", err)
		}
		cfg.PublicKey = publicKey
	}

	return &cfg, nil
}

func (c *AgentConfig) init() {
	flag.StringVar(&c.ServerAddress, "a", serverAddressDefault, "Start server address (default - :8080)")
	flag.IntVar(&c.ReportInterval, "r", reportIntervalDefault, "Interval of report metric")
	flag.IntVar(&c.PollInterval, "p", pollIntervalDefault, "Interval of poll metric")
	flag.StringVar(&c.SignKey, "k", "", "Server key")
	flag.IntVar(&c.RateLimit, "l", rateLimitDefault, "Rate limit")
	flag.StringVar(&c.PublicKeyPath, "-crypto-key", "", "Public key path")
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

func (c *AgentConfig) getPublicKey() (*rsa.PublicKey, error) {
	publicKeyPEM, err := os.ReadFile(c.PublicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	if publicKeyBlock == nil {
		return nil, err
	}

	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}
