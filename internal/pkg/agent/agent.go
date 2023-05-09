package agent

import (
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent/config"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartClient(c *config.AgentConfig) {

	pollInterval := time.Duration(c.PollInterval) * time.Second     //тесты не проходят с duration
	reportInterval := time.Duration(c.ReportInterval) * time.Second //тесты не проходят с duration

	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	m := storage.NewMetrics()

	go func() {
		for {
			select {
			case <-pollTicker.C:
				if err := m.UpdateMetrics(); err != nil {
					log.Fatal(err)
				}
			case <-reportTicker.C:
				if err := m.SendMetrics(c.ServerAddress); err != nil {
					log.Fatal(err)
				}
			}
		}
	}()

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChanel
}
