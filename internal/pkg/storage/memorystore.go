package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Memory interface {
	UpdateGaugeMetric(name string, value metrics.Gauge) error
	UpdateCounterMetric(name string, value metrics.Counter) error
}

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

func NewMetricsFile(file string, storeInterval time.Duration) *MemoryStore {
	return &MemoryStore{
		metrics:         make(map[string]*metrics.Metrics),
		fileStoragePath: file,
		storeInterval:   storeInterval,
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

func (m *MemoryStore) UpdateCounterMetric(metricName string, metricValue metrics.Counter) error {
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

func (m *MemoryStore) Close() error {
	if m.tickerDone != nil {
		m.tickerDone <- struct{}{}
	}
	return nil
}

func UpdateMetrics(m Memory) error {
	var metricsStats runtime.MemStats
	runtime.ReadMemStats(&metricsStats)
	var errorsSlice []string

	err := m.UpdateGaugeMetric("Alloc", metrics.Gauge(metricsStats.Alloc))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("BuckHashSys", metrics.Gauge(metricsStats.BuckHashSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("BuckHashSys", metrics.Gauge(metricsStats.BuckHashSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("Frees", metrics.Gauge(metricsStats.Frees))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("GCCPUFraction", metrics.Gauge(metricsStats.GCCPUFraction))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("GCSys", metrics.Gauge(metricsStats.GCSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("HeapAlloc", metrics.Gauge(metricsStats.HeapAlloc))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("HeapIdle", metrics.Gauge(metricsStats.HeapIdle))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("HeapInuse", metrics.Gauge(metricsStats.HeapInuse))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("HeapObjects", metrics.Gauge(metricsStats.HeapObjects))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("HeapReleased", metrics.Gauge(metricsStats.HeapReleased))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("HeapSys", metrics.Gauge(metricsStats.HeapSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("LastGC", metrics.Gauge(metricsStats.LastGC))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("Lookups", metrics.Gauge(metricsStats.Lookups))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("MCacheInuse", metrics.Gauge(metricsStats.MCacheInuse))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("MCacheSys", metrics.Gauge(metricsStats.MCacheSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("MSpanInuse", metrics.Gauge(metricsStats.MSpanInuse))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("MSpanSys", metrics.Gauge(metricsStats.MSpanSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("Mallocs", metrics.Gauge(metricsStats.Mallocs))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("NextGC", metrics.Gauge(metricsStats.NextGC))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("NumForcedGC", metrics.Gauge(metricsStats.NumForcedGC))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("NumGC", metrics.Gauge(metricsStats.NumGC))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("OtherSys", metrics.Gauge(metricsStats.OtherSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("PauseTotalNs", metrics.Gauge(metricsStats.PauseTotalNs))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("StackInuse", metrics.Gauge(metricsStats.StackInuse))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("StackSys", metrics.Gauge(metricsStats.StackSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("Sys", metrics.Gauge(metricsStats.Sys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("TotalAlloc", metrics.Gauge(metricsStats.TotalAlloc))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric("RandomValue", metrics.Gauge(rand.Float64()))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}

	err = m.UpdateCounterMetric("PollCount", 1)
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}

	if len(errorsSlice) > 0 {
		return fmt.Errorf(strings.Join(errorsSlice, "\n"))
	} else {
		return nil
	}
}
