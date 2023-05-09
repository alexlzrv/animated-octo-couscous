package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	ServerAddress string `env:"ADDRESS"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func Init(c *ServerConfig) {
	flag.StringVar(&c.ServerAddress, "a", "localhost:8080", "Listen server address (default - :8080)")
	flag.Parse()

	if envServerAddress := os.Getenv("ADDRESS"); envServerAddress != "" {
		c.ServerAddress = envServerAddress
	}
}
