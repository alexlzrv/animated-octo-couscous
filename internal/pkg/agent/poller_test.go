package agent_test

import (
	"context"
	"testing"

	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
)

func TestPoolWorker(t *testing.T) {
	mtr := storage.NewMetrics()
	agent.UpdateMetrics(context.Background(), mtr)

	counterMetric, _ := mtr.GetMetric(context.Background(), "PollCount", "")
	if *counterMetric.Delta != 1 {
		t.Errorf("Counter wasn't incremented: %d", *counterMetric.Delta)
	}
}
