package storage

import (
	"context"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
)

type Store interface {
	UpdateCounterMetric(ctx context.Context, name string, value metrics.Counter) error
	UpdateGaugeMetric(ctx context.Context, name string, value metrics.Gauge) error
	UpdateMetrics(ctx context.Context, metricBatch []*metrics.Metrics) error

	ResetCounterMetric(ctx context.Context, name string) error

	GetMetric(ctx context.Context, name string, metricType string) (*metrics.Metrics, bool)
	GetMetrics(ctx context.Context) (map[string]*metrics.Metrics, error)

	LoadMetrics(filePath string) error
	SaveMetrics(filePath string) error

	Ping() error
	Close() error
}
