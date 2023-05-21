package storage

import (
	"fmt"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"math/rand"
	"runtime"
	"sync"
)

type MemoryStore struct {
	metrics map[string]*metrics.Metrics
	lock    sync.Mutex
}

func NewMetrics() *MemoryStore {
	return &MemoryStore{
		metrics: make(map[string]*metrics.Metrics),
	}
}

func (m *MemoryStore) UpdateGaugeMetric(metricName string, metricValue metrics.Gauge) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	currentMetric, ok := m.metrics[metricName]

	switch {
	case ok && currentMetric.Value != nil:
		*(currentMetric.Value) = metricValue
	case ok && currentMetric.Value == nil:
		return fmt.Errorf("%s %s", metricName, currentMetric.MType)
	default:
		m.metrics[metricName] = &metrics.Metrics{
			ID:    metricName,
			MType: metrics.GaugeMetricName,
			Value: &metricValue,
		}
	}
	return nil
}

func (m *MemoryStore) UpdateCounterMetric(metricName string, metricValue metrics.Counter) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	currentMetric, ok := m.metrics[metricName]

	switch {
	case ok && currentMetric.Delta != nil:
		*(currentMetric.Delta) += metricValue
	case ok && currentMetric.Delta == nil:
		return fmt.Errorf("%s %s", metricName, currentMetric.MType)
	default:
		m.metrics[metricName] = &metrics.Metrics{
			ID:    metricName,
			MType: metrics.CounterMetricName,
			Delta: &metricValue,
		}
	}

	return nil
}

func (m *MemoryStore) GetMetric(metricName string) (*metrics.Metrics, bool) {
	metric, ok := m.metrics[metricName]

	return metric, ok
}

func (m *MemoryStore) GetMetrics() map[string]*metrics.Metrics {
	return m.metrics
}

func (m *MemoryStore) ResetCounterMetric(metricName string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	var zero metrics.Counter
	currentMetric, ok := m.metrics[metricName]
	switch {
	case ok && currentMetric.Delta != nil:
		*(currentMetric.Delta) = zero
	case ok && currentMetric.Delta == nil:
		return fmt.Errorf("%s %s", metricName, currentMetric.MType)
	default:
		m.metrics[metricName] = &metrics.Metrics{
			ID:    metricName,
			MType: metrics.CounterMetricName,
			Delta: &zero,
		}
	}
	return nil
}

func UpdateMetrics(m Storage) error {
	var metricsStats runtime.MemStats
	runtime.ReadMemStats(&metricsStats)

	_ = m.UpdateGaugeMetric("Alloc", metrics.Gauge(metricsStats.Alloc))
	_ = m.UpdateGaugeMetric("BuckHashSys", metrics.Gauge(metricsStats.BuckHashSys))
	_ = m.UpdateGaugeMetric("BuckHashSys", metrics.Gauge(metricsStats.BuckHashSys))
	_ = m.UpdateGaugeMetric("Frees", metrics.Gauge(metricsStats.Frees))
	_ = m.UpdateGaugeMetric("GCCPUFraction", metrics.Gauge(metricsStats.GCCPUFraction))
	_ = m.UpdateGaugeMetric("GCSys", metrics.Gauge(metricsStats.GCSys))
	_ = m.UpdateGaugeMetric("HeapAlloc", metrics.Gauge(metricsStats.HeapAlloc))
	_ = m.UpdateGaugeMetric("HeapIdle", metrics.Gauge(metricsStats.HeapIdle))
	_ = m.UpdateGaugeMetric("HeapInuse", metrics.Gauge(metricsStats.HeapInuse))
	_ = m.UpdateGaugeMetric("HeapObjects", metrics.Gauge(metricsStats.HeapObjects))
	_ = m.UpdateGaugeMetric("HeapReleased", metrics.Gauge(metricsStats.HeapReleased))
	_ = m.UpdateGaugeMetric("HeapSys", metrics.Gauge(metricsStats.HeapSys))
	_ = m.UpdateGaugeMetric("LastGC", metrics.Gauge(metricsStats.LastGC))
	_ = m.UpdateGaugeMetric("Lookups", metrics.Gauge(metricsStats.Lookups))
	_ = m.UpdateGaugeMetric("MCacheInuse", metrics.Gauge(metricsStats.MCacheInuse))
	_ = m.UpdateGaugeMetric("MCacheSys", metrics.Gauge(metricsStats.MCacheSys))
	_ = m.UpdateGaugeMetric("MSpanInuse", metrics.Gauge(metricsStats.MSpanInuse))
	_ = m.UpdateGaugeMetric("MSpanSys", metrics.Gauge(metricsStats.MSpanSys))
	_ = m.UpdateGaugeMetric("Mallocs", metrics.Gauge(metricsStats.Mallocs))
	_ = m.UpdateGaugeMetric("NextGC", metrics.Gauge(metricsStats.NextGC))
	_ = m.UpdateGaugeMetric("NumForcedGC", metrics.Gauge(metricsStats.NumForcedGC))
	_ = m.UpdateGaugeMetric("NumGC", metrics.Gauge(metricsStats.NumGC))
	_ = m.UpdateGaugeMetric("OtherSys", metrics.Gauge(metricsStats.OtherSys))
	_ = m.UpdateGaugeMetric("PauseTotalNs", metrics.Gauge(metricsStats.PauseTotalNs))
	_ = m.UpdateGaugeMetric("StackInuse", metrics.Gauge(metricsStats.StackInuse))
	_ = m.UpdateGaugeMetric("StackSys", metrics.Gauge(metricsStats.StackSys))
	_ = m.UpdateGaugeMetric("Sys", metrics.Gauge(metricsStats.Sys))
	_ = m.UpdateGaugeMetric("TotalAlloc", metrics.Gauge(metricsStats.TotalAlloc))
	_ = m.UpdateGaugeMetric("RandomValue", metrics.Gauge(rand.Float64()))

	_ = m.UpdateCounterMetric("PollCount", 1)

	return nil
}
