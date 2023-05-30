package storage

import "github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"

type Storage interface {
	UpdateCounterMetric(name string, value metrics.Counter) error
	ResetCounterMetric(name string) error
	UpdateGaugeMetric(name string, value metrics.Gauge) error
	GetMetric(name string) (*metrics.Metrics, bool)
	GetMetrics() map[string]*metrics.Metrics
	LoadMetrics(filePath string) error
	SaveMetrics(filePath string) error
}
