package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	ServerAddress   string `env:"ADDRESS" json:"server_address"`
	SignKey         string `env:"KEY"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	StoreInterval   int    `env:"STORE_INTERVAL" json:"store_interval"`
	Restore         bool   `env:"RESTORE" json:"restore"`
	PrivateKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigPath      string `env:"CONFIG"`
	SignKeyByte     []byte
}

const (
	storeIntervalDefault = 300
	serverAddressDefault = "localhost:8080"
	filePathDefault      = "/tmp/metrics-db.json"
)

func NewServerConfig() (*ServerConfig, error) {
	cfg := ServerConfig{}
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

func (c *ServerConfig) init() {
	flag.StringVar(&c.ServerAddress, "a", serverAddressDefault, "Listen server address (default - :8080)")
	flag.IntVar(&c.StoreInterval, "i", storeIntervalDefault, "Store interval")
	flag.StringVar(&c.FileStoragePath, "f", filePathDefault, "File storage path")
	flag.BoolVar(&c.Restore, "r", true, "Restore")
	flag.StringVar(&c.DatabaseDSN, "d", "", "Connect database string")
	flag.StringVar(&c.SignKey, "k", "", "Server key")
	flag.StringVar(&c.PrivateKey, "-crypto-key", "", "Private key path")
	flag.StringVar(&c.ConfigPath, "c", "", "Path to config file")
	flag.StringVar(&c.ConfigPath, "config", "", "Path to config file (the same as -c)")
	flag.Parse()
}

func readConfigFile(path string) (cfg *ServerConfig, err error) {
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
