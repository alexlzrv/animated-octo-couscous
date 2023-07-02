package agent

import (
	"context"
	"sync"
	"time"

	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent/config"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
)

func StartClient(ctx context.Context, c *config.AgentConfig) {
	logrus.Info("Agent is running...")
	wg := &sync.WaitGroup{}

	metric := storage.NewMetrics()

	pollerInterval := time.Duration(c.PollInterval) * time.Second
	pollerTicker := time.NewTicker(pollerInterval)
	defer pollerTicker.Stop()

	reportInterval := time.Duration(c.ReportInterval) * time.Second
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		RunUpdateMemStatMetrics(ctx, pollerTicker, metric)
	}()

	for i := 1; i < c.RateLimit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			RunSendMetric(ctx, reportTicker, c, metric)
		}()
	}

	wg.Wait()
}
