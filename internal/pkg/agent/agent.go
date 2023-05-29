package agent

import (
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent/config"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartClient(c *config.AgentConfig) {
	logrus.Info("Agent is running...")
	pollInterval := time.Duration(c.PollInterval) * time.Second     //тесты не проходят с duration
	reportInterval := time.Duration(c.ReportInterval) * time.Second //тесты не проходят с duration

	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	metric := storage.NewMetrics()

	go func() {
		for {
			select {
			case <-pollTicker.C:
				if err := storage.UpdateMetrics(metric); err != nil {
					logrus.Errorf("Error update metrics %s", err)
				}
			case <-reportTicker.C:
				if err := SendMetrics(metric, c.ServerAddress); err != nil {
					logrus.Errorf("Error send metrics %s", err)
				}
				if err := metric.ResetCounterMetric("PollCount"); err != nil {
					logrus.Errorf("Error reset metrics %s", err)
				}
			}
		}
	}()

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChanel
}
