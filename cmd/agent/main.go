package main

import (
	"context"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent/config"
)

func main() {
	cfg := config.NewAgentConfig()
	agent.StartClient(context.Background(), cfg)
}
