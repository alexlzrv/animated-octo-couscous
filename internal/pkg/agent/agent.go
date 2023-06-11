package agent

import (
	"context"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent/config"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartClient(ctx context.Context, c *config.AgentConfig) {
	logrus.Info("Agent is running...")
	pollInterval := time.Duration(c.PollInterval) * time.Second     //тесты не проходят с duration
	reportInterval := time.Duration(c.ReportInterval) * time.Second //тесты не проходят с duration

	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	metric := storage.NewMetrics()

	agentContext, cancelCtx := context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-pollTicker.C:
				if err := UpdateMetrics(agentContext, metric); err != nil {
					logrus.Errorf("Error update metrics %v", err)
				}
			case <-reportTicker.C:
				if err := SendMetrics(agentContext, metric, c.ServerAddress); err != nil {
					logrus.Errorf("Error send metrics %v", err)
				}
				if err := metric.ResetCounterMetric(agentContext, "PollCount"); err != nil {
					logrus.Errorf("Error reset metrics %v", err)
				}
			}
		}
	}()
	cancelCtx()
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChanel
}
