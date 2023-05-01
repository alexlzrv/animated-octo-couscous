package storage

const (
	GaugeMetricName   = "gauge"
	CounterMetricName = "counter"
)

type Gauge float64
type Counter int64

type Metrics struct {
	GaugeMetrics   map[string]Gauge
	CounterMetrics map[string]Counter
}

func NewMetrics() *Metrics {
	var metric Metrics
	metric.GaugeMetrics = make(map[string]Gauge)
	metric.CounterMetrics = make(map[string]Counter)
	return &metric
}
