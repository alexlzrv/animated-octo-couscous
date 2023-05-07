package main

import (
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/server"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	metricServer := server.MetricsServer{
		MetricsStore: storage.NewMetrics(),
	}

	go metricServer.StartListener()

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChanel
}
