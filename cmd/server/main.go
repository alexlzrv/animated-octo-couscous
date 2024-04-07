package main

import (
	"context"
	"log"

	"github.com/mayr0y/animated-octo-couscous.git/internal/greetings"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server/config"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	if err := greetings.Hello(buildVersion, buildDate, buildCommit); err != nil {
		log.Fatal(err)
	}

	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatal(err)
	}
	server.StartListener(context.Background(), cfg)
}
