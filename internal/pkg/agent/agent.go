package agent

import (
	"fmt"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
)

type gauge float64

type Config struct {
	ServerAddress  string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func NewAgentConfig() *Config {
	return &Config{}
}

func GetMetrics() map[string]gauge {
	var metrics runtime.MemStats
	runtime.ReadMemStats(&metrics)

	metricsMap := map[string]gauge{
		"Alloc":         gauge(metrics.Alloc),
		"BuckHashSys":   gauge(metrics.BuckHashSys),
		"Frees":         gauge(metrics.Frees),
		"GCCPUFraction": gauge(metrics.GCCPUFraction),
		"GCSys":         gauge(metrics.GCSys),
		"HeapAlloc":     gauge(metrics.HeapAlloc),
		"HeapIdle":      gauge(metrics.HeapIdle),
		"HeapInuse":     gauge(metrics.HeapInuse),
		"HeapObjects":   gauge(metrics.HeapObjects),
		"HeapReleased":  gauge(metrics.HeapReleased),
		"HeapSys":       gauge(metrics.HeapSys),
		"LastGC":        gauge(metrics.LastGC),
		"Lookups":       gauge(metrics.Lookups),
		"MCacheInuse":   gauge(metrics.MCacheInuse),
		"MCacheSys":     gauge(metrics.MCacheSys),
		"MSpanInuse":    gauge(metrics.MSpanInuse),
		"MSpanSys":      gauge(metrics.MSpanSys),
		"Mallocs":       gauge(metrics.Mallocs),
		"NextGC":        gauge(metrics.NextGC),
		"NumForcedGC":   gauge(metrics.NumForcedGC),
		"NumGC":         gauge(metrics.NumGC),
		"OtherSys":      gauge(metrics.OtherSys),
		"PauseTotalNs":  gauge(metrics.PauseTotalNs),
		"StackInuse":    gauge(metrics.StackInuse),
		"StackSys":      gauge(metrics.StackSys),
		"Sys":           gauge(metrics.Sys),
		"TotalAlloc":    gauge(metrics.TotalAlloc),
		"RandomValue":   gauge(rand.Float64()),
	}

	return metricsMap
}

func SendMetrics(metric map[string]gauge, count int64, serverAddress string) error {
	for name, value := range metric {
		urlCounter := createURL(serverAddress, storage.CounterMetricName, "PollCount", fmt.Sprintf("%d", count))
		err := createPostRequest(urlCounter)
		if err != nil {
			return err
		}

		urlGauge := createURL(serverAddress, storage.GaugeMetricName, name, fmt.Sprintf("%f", float64(value)))
		err = createPostRequest(urlGauge)
		if err != nil {
			return err
		}
	}
	return nil
}

func createURL(address string, metricType string, metricName string, metricValue string) string {
	str := []string{"http:/", address, "update", metricType, metricName, metricValue}
	return strings.Join(str, "/")
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
