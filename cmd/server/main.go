package main

import (
	"flag"
	"github.com/caarlos0/env"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	metricServer := server.MetricsServer{
		MetricsStore: storage.NewMetrics(),
	}

	config := server.NewServerConfig()
	flag.StringVar(&config.ServerAddress, "a", "localhost:8080", "ServerAddress")
	flag.Parse()

	if err := env.Parse(config); err != nil {
		log.Fatalf(err.Error())
	}

	go metricServer.StartListener(config)

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChanel
}
