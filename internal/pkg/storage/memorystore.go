package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
)

type MemoryStore struct {
	Metrics         map[string]*metrics.Metrics
	FileStoragePath string
	storeInterval   time.Duration
	tickerDone      chan struct{}
	lock            sync.Mutex
	db              *sql.DB
}

func NewMetrics() *MemoryStore {
	return &MemoryStore{
		Metrics: make(map[string]*metrics.Metrics),
	}
}

func NewMetricsFile(file string, storeInterval time.Duration) (*MemoryStore, error) {
	metricStore := MemoryStore{
		Metrics:         make(map[string]*metrics.Metrics),
		FileStoragePath: file,
		storeInterval:   storeInterval,
	}

	return &metricStore, nil
}

func (m *MemoryStore) UpdateMetrics(_ context.Context, metricBatch []*metrics.Metrics) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, metric := range metricBatch {
		currentMetric, ok := m.Metrics[metric.ID]
		switch {
		case ok && metric.MType == metrics.GaugeMetricName && currentMetric.Value != nil:
			currentMetric.Value = metric.Value
		case ok && metric.MType == metrics.GaugeMetricName && currentMetric.Value == nil:
			return fmt.Errorf("mismatch metric type %s:%s", metric.ID, currentMetric.MType)
		case ok && metric.MType == metrics.CounterMetricName && currentMetric.Delta != nil:
			*(currentMetric.Delta) += *(metric.Delta)
		case ok && metric.MType == metrics.CounterMetricName && currentMetric.Delta == nil:
			return fmt.Errorf("mismatch metric type %s:%s", metric.ID, currentMetric.MType)
		default:
			m.Metrics[metric.ID] = metric
		}
	}

	return nil
}

func (m *MemoryStore) UpdateGaugeMetric(_ context.Context, metricName string, metricValue metrics.Gauge) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	currentMetric, ok := m.Metrics[metricName]

	switch {
	case ok && currentMetric.Value != nil:
		*(currentMetric.Value) = metricValue
	case ok && currentMetric.Value == nil:
		return fmt.Errorf("mismatch metric type %s:%s", metricName, currentMetric.MType)
	default:
		m.Metrics[metricName] = &metrics.Metrics{
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
	currentMetric, ok := m.Metrics[metricName]

	switch {
	case ok && currentMetric.Delta != nil:
		*(currentMetric.Delta) += metricValue
	case ok && currentMetric.Delta == nil:
		return fmt.Errorf("mismatch metric type %s:%s", metricName, currentMetric.MType)
	default:
		m.Metrics[metricName] = &metrics.Metrics{
			ID:    metricName,
			MType: metrics.CounterMetricName,
			Delta: &metricValue,
		}
	}

	return nil
}

func (m *MemoryStore) GetMetric(_ context.Context, metricName string, _ string) (*metrics.Metrics, bool) {
	metric, ok := m.Metrics[metricName]

	return metric, ok
}

func (m *MemoryStore) GetMetrics(_ context.Context) (map[string]*metrics.Metrics, error) {
	return m.Metrics, nil
}

func (m *MemoryStore) ResetCounterMetric(_ context.Context, metricName string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	var zero metrics.Counter
	currentMetric, ok := m.Metrics[metricName]
	switch {
	case ok && currentMetric.Delta != nil:
		*(currentMetric.Delta) = zero
	case ok && currentMetric.Delta == nil:
		return fmt.Errorf("mismatch metric type %s:%s", metricName, currentMetric.MType)
	default:
		m.Metrics[metricName] = &metrics.Metrics{
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
	return jsonDecoder.Decode(&m.Metrics)
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
	return encoder.Encode(&m.Metrics)
}

func (m *MemoryStore) Ping() error {
	return nil
}

func (m *MemoryStore) Close() error {
	return nil
}
