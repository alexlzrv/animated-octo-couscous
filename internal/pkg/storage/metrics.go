package storage

import (
	"math/rand"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"sync"
)

type Metric interface {
	UpdateGaugeMetric(name string, value Gauge)
	UpdateCounterMetric(name string, value Counter)
	GetGaugeMetric(name string) (Gauge, bool)
	GetCounterMetric(name string) (Counter, bool)
	GetGaugeMetrics() map[string]Gauge
	GetCounterMetrics() map[string]Counter
}

const (
	GaugeMetricName   = "gauge"
	CounterMetricName = "counter"
	pollCountName     = "PollCount"
)

type Gauge float64
type Counter int64

type Metrics struct {
	GaugeMetrics   map[string]Gauge
	CounterMetrics map[string]Counter
	lock           sync.Mutex
}

func NewMetrics() *Metrics {
	m := Metrics{
		GaugeMetrics:   make(map[string]Gauge),
		CounterMetrics: make(map[string]Counter),
	}
	m.CounterMetrics[pollCountName] = 0
	return &m
}

func (m *Metrics) UpdateGaugeMetric(metricName string, metricValue Gauge) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.GaugeMetrics[metricName] = metricValue
}

func (m *Metrics) UpdateCounterMetric(metricName string, metricValue Counter) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.CounterMetrics[metricName] += metricValue
}

func (m *Metrics) GetGaugeMetric(metricName string) (Gauge, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	metric, ok := m.GaugeMetrics[metricName]

	return metric, ok
}

func (m *Metrics) GetCounterMetric(metricName string) (Counter, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	metric, ok := m.CounterMetrics[metricName]

	return metric, ok
}

func (m *Metrics) GetGaugeMetrics() map[string]Gauge {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.GaugeMetrics
}

func (m *Metrics) GetCounterMetrics() map[string]Counter {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.CounterMetrics
}

func (m *Metrics) UpdateMetrics() error {
	var metrics runtime.MemStats
	runtime.ReadMemStats(&metrics)

	m.GaugeMetrics["Alloc"] = Gauge(metrics.Alloc)
	m.GaugeMetrics["BuckHashSys"] = Gauge(metrics.BuckHashSys)
	m.GaugeMetrics["Frees"] = Gauge(metrics.Frees)
	m.GaugeMetrics["GCCPUFraction"] = Gauge(metrics.GCCPUFraction)
	m.GaugeMetrics["GCSys"] = Gauge(metrics.GCSys)
	m.GaugeMetrics["HeapAlloc"] = Gauge(metrics.HeapAlloc)
	m.GaugeMetrics["HeapIdle"] = Gauge(metrics.HeapIdle)
	m.GaugeMetrics["HeapInuse"] = Gauge(metrics.HeapInuse)
	m.GaugeMetrics["HeapObjects"] = Gauge(metrics.HeapObjects)
	m.GaugeMetrics["HeapReleased"] = Gauge(metrics.HeapReleased)
	m.GaugeMetrics["HeapSys"] = Gauge(metrics.HeapSys)
	m.GaugeMetrics["LastGC"] = Gauge(metrics.LastGC)
	m.GaugeMetrics["Lookups"] = Gauge(metrics.Lookups)
	m.GaugeMetrics["MCacheInuse"] = Gauge(metrics.MCacheInuse)
	m.GaugeMetrics["MCacheSys"] = Gauge(metrics.MCacheSys)
	m.GaugeMetrics["MSpanInuse"] = Gauge(metrics.MSpanInuse)
	m.GaugeMetrics["MSpanSys"] = Gauge(metrics.MSpanSys)
	m.GaugeMetrics["Mallocs"] = Gauge(metrics.Mallocs)
	m.GaugeMetrics["NextGC"] = Gauge(metrics.NextGC)
	m.GaugeMetrics["NumForcedGC"] = Gauge(metrics.NumForcedGC)
	m.GaugeMetrics["NumGC"] = Gauge(metrics.NumGC)
	m.GaugeMetrics["OtherSys"] = Gauge(metrics.OtherSys)
	m.GaugeMetrics["PauseTotalNs"] = Gauge(metrics.PauseTotalNs)
	m.GaugeMetrics["StackInuse"] = Gauge(metrics.StackInuse)
	m.GaugeMetrics["StackSys"] = Gauge(metrics.StackSys)
	m.GaugeMetrics["Sys"] = Gauge(metrics.Sys)
	m.GaugeMetrics["TotalAlloc"] = Gauge(metrics.TotalAlloc)
	m.GaugeMetrics["RandomValue"] = Gauge(rand.Float64())

	m.CounterMetrics[pollCountName]++

	return nil
}

func (m *Metrics) SendMetrics(serverAddress string) error {
	for name, value := range m.GaugeMetrics {
		urlGauge, err := createURL(serverAddress, GaugeMetricName, name, strconv.FormatFloat(float64(value), 'E', -1, 32))
		if err != nil {
			return err
		}
		err = createPostRequest(urlGauge)
		if err != nil {
			return err
		}
	}

	for name, value := range m.CounterMetrics {
		urlGauge, err := createURL(serverAddress, CounterMetricName, name, strconv.FormatFloat(float64(value), 'E', -1, 32))
		if err != nil {
			return err
		}
		err = createPostRequest(urlGauge)
		if err != nil {
			return err
		}
		if name == pollCountName {
			m.CounterMetrics[pollCountName] = 0
		}
	}

	return nil
}

func createURL(address string, metricType string, metricName string, metricValue string) (string, error) {
	return url.JoinPath("https://", address, "update", metricType, metricName, metricValue)
}

func createPostRequest(url string) error {
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
