package storage

import (
	"context"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
)

type Store interface {
	UpdateCounterMetric(ctx context.Context, name string, value metrics.Counter) error
	ResetCounterMetric(ctx context.Context, name string) error
	UpdateGaugeMetric(ctx context.Context, name string, value metrics.Gauge) error
	GetMetric(ctx context.Context, name string, metricType string) (*metrics.Metrics, bool)
	GetMetrics(ctx context.Context) (map[string]*metrics.Metrics, error)
	Ping(ctx context.Context) error
	LoadMetrics(filePath string) error
	SaveMetrics(filePath string) error
}
