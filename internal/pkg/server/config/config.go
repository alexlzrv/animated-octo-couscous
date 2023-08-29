package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	ServerAddress   string `env:"ADDRESS"`
	SignKey         string `env:"KEY"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
	SignKeyByte     []byte
}

const (
	storeIntervalDefault = 300
	serverAddressDefault = "localhost:8080"
	filePathDefault      = "/tmp/metrics-db.json"
)

func NewServerConfig() *ServerConfig {
	cfg := ServerConfig{}
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

func (c *ServerConfig) init() {
	flag.StringVar(&c.ServerAddress, "a", serverAddressDefault, "Listen server address (default - :8080)")
	flag.IntVar(&c.StoreInterval, "i", storeIntervalDefault, "Store interval")
	flag.StringVar(&c.FileStoragePath, "f", filePathDefault, "File storage path")
	flag.BoolVar(&c.Restore, "r", true, "Restore")
	flag.StringVar(&c.DatabaseDSN, "d", "", "Connect database string")
	flag.StringVar(&c.SignKey, "k", "", "Server key")
	flag.Parse()
}
