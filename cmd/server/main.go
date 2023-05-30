package main

import (
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server/config"
)

func main() {
	cfg := config.NewServerConfig()
	server.StartListener(cfg)
}
