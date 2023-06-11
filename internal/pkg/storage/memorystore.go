package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"os"
	"sync"
	"time"
)

type MemoryStore struct {
	metrics         map[string]*metrics.Metrics
	fileStoragePath string
	storeInterval   time.Duration
	tickerDone      chan struct{}
	lock            sync.Mutex
	db              *sql.DB
}

func NewMetrics() *MemoryStore {
	return &MemoryStore{
		metrics: make(map[string]*metrics.Metrics),
	}
}

func NewMetricsFile(file string, storeInterval time.Duration) (*MemoryStore, error) {
	metricStore := MemoryStore{
		metrics:         make(map[string]*metrics.Metrics),
		fileStoragePath: file,
		storeInterval:   storeInterval,
	}

	return &metricStore, nil
}

func (m *MemoryStore) UpdateGaugeMetric(_ context.Context, metricName string, metricValue metrics.Gauge) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	currentMetric, ok := m.metrics[metricName]

	switch {
	case ok && currentMetric.Value != nil:
		*(currentMetric.Value) = metricValue
	case ok && currentMetric.Value == nil:
		return fmt.Errorf("mismatch metric type %s:%s", metricName, currentMetric.MType)
	default:
		m.metrics[metricName] = &metrics.Metrics{
			ID:    metricName,
			MType: metrics.GaugeMetricName,
			Value: &metricValue,
		}
	}
	return nil
}

func (m *MemoryStore) UpdateCounterMetric(_ context.Context, metricName string, metricValue metrics.Counter) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	currentMetric, ok := m.metrics[metricName]

	switch {
	case ok && currentMetric.Delta != nil:
		*(currentMetric.Delta) += metricValue
	case ok && currentMetric.Delta == nil:
		return fmt.Errorf("mismatch metric type %s:%s", metricName, currentMetric.MType)
	default:
		m.metrics[metricName] = &metrics.Metrics{
			ID:    metricName,
			MType: metrics.CounterMetricName,
			Delta: &metricValue,
		}
	}

	return nil
}

func (m *MemoryStore) GetMetric(_ context.Context, metricName string, _ string) (*metrics.Metrics, bool) {
	metric, ok := m.metrics[metricName]

	return metric, ok
}

func (m *MemoryStore) GetMetrics(_ context.Context) (map[string]*metrics.Metrics, error) {
	return m.metrics, nil
}

func (m *MemoryStore) ResetCounterMetric(_ context.Context, metricName string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	var zero metrics.Counter
	currentMetric, ok := m.metrics[metricName]
	switch {
	case ok && currentMetric.Delta != nil:
		*(currentMetric.Delta) = zero
	case ok && currentMetric.Delta == nil:
		return fmt.Errorf("mismatch metric type %s:%s", metricName, currentMetric.MType)
	default:
		m.metrics[metricName] = &metrics.Metrics{
			ID:    metricName,
			MType: metrics.CounterMetricName,
			Delta: &zero,
		}
	}
	return nil
}

func (m *MemoryStore) LoadMetrics(filePath string) error {
	if filePath == "" {
		return nil
	}
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	jsonDecoder := json.NewDecoder(file)
	return jsonDecoder.Decode(&m.metrics)
}

func (m *MemoryStore) SaveMetrics(filePath string) error {
	if filePath == "" {
		return nil
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(&m.metrics)
}

func (m *MemoryStore) Ping() error {
	return nil
}

func (m *MemoryStore) Close() error {
	return nil
}
