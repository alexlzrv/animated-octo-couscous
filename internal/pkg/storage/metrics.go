package storage

const (
	GaugeMetricName   = "gauge"
	CounterMetricName = "counter"
)

type Gauge float64
type Counter int64

type Metric interface {
	UpdateGaugeMetric(name string, value Gauge)
	UpdateCounterMetric(name string, value Counter)
	GetGaugeMetric(name string) (Gauge, bool)
	GetCounterMetric(name string) (Counter, bool)
	GetGaugeMetrics() map[string]Gauge
	GetCounterMetrics() map[string]Counter
}

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

func (m *Metrics) UpdateGaugeMetric(metricName string, metricValue Gauge) {
	m.GaugeMetrics[metricName] = metricValue
}

func (m *Metrics) UpdateCounterMetric(metricName string, metricValue Counter) {
	m.CounterMetrics[metricName] += metricValue
}

func (m *Metrics) GetGaugeMetric(metricName string) (Gauge, bool) {
	metric, ok := m.GaugeMetrics[metricName]

	return metric, ok
}

func (m *Metrics) GetCounterMetric(metricName string) (Counter, bool) {
	metric, ok := m.CounterMetrics[metricName]

	return metric, ok
}

func (m *Metrics) GetGaugeMetrics() map[string]Gauge {
	return m.GaugeMetrics
}

func (m *Metrics) GetCounterMetrics() map[string]Counter {
	return m.CounterMetrics
}
