package main

import (
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent/config"
)

func main() {
	cfg := config.NewAgentConfig()
	config.Init(cfg)
	agent.StartClient(cfg)
}
