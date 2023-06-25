package agent

import (
	"context"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent/config"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
	"os/signal"
	"sync"
	"syscall"
)

func StartClient(ctx context.Context, c *config.AgentConfig) {
	logrus.Info("Agent is running...")
	ctx, stop := signal.NotifyContext(ctx,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()
	wg := sync.WaitGroup{}

	metric := storage.NewMetrics()

	wg.Add(1)
	go func() {
		defer wg.Done()
		RunUpdateMemStatMetrics(ctx, c, metric)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		RunUpdateVirtualMetrics(ctx, c, metric)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		RunSendMetric(ctx, c, metric)
	}()

	wg.Wait()
}
