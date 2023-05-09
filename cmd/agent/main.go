package main

import (
	"flag"
	"github.com/caarlos0/env"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := agent.NewAgentConfig()

	flag.StringVar(&config.ServerAddress, "a", "localhost:8080", "ServerAddress")
	flag.IntVar(&config.ReportInterval, "r", 10, "ReportInterval")
	flag.IntVar(&config.PollInterval, "p", 2, "PollInterval")
	flag.Parse()

	if err := env.Parse(config); err != nil {
		log.Fatalf(err.Error())
	}

	var pollCount int64

	pollInterval := time.Duration(config.PollInterval) * time.Second
	reportInterval := time.Duration(config.ReportInterval) * time.Second

	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	go func() {
		metric := agent.GetMetrics()
		for {
			select {
			case <-pollTicker.C:
				pollCount++
				metric = agent.GetMetrics()
			case <-reportTicker.C:
				if err := agent.SendMetrics(metric, pollCount, config.ServerAddress); err != nil {
					log.Println(err.Error())
				}
				pollCount = 0
			}
		}
	}()

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChanel
}
