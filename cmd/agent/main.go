package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/mayr0y/animated-octo-couscous.git/internal/greetings"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent/config"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	if err := greetings.Hello(buildVersion, buildDate, buildCommit); err != nil {
		log.Fatal(err)
	}

	cfg, err := config.NewAgentConfig()
	if err != nil {
		log.Fatal(err)
	}
	agent.StartClient(ctx, cfg)
}
