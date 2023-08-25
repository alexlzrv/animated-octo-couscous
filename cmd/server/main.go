package main

import (
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

	cfg := config.NewServerConfig()
	server.StartListener(cfg)
}
