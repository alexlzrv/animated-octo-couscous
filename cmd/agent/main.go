package main

import (
	"fmt"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type gauge float64

func getMetrics() map[string]gauge {
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

func main() {
	var pollCount int64

	pollInterval := time.NewTicker(2 * time.Second)
	defer pollInterval.Stop()

	reportInterval := time.NewTicker(10 * time.Second)
	defer reportInterval.Stop()

	go func() {
		metric := getMetrics()
		for {
			select {
			case <-pollInterval.C:
				pollCount++
				metric = getMetrics()
			case <-reportInterval.C:
				if err := sendMetrics(metric, pollCount); err != nil {
					log.Println(err.Error())
				}
				pollCount = 0
			}
		}
	}()

	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
	<-signalChanel
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

func sendMetrics(metric map[string]gauge, count int64) error {
	for name, value := range metric {
		urlCounter := createURL("localhost:8080", storage.CounterMetricName, "PollCount", fmt.Sprintf("%d", count))
		err := createPostRequest(urlCounter)
		if err != nil {
			return err
		}

		urlGauge := createURL("localhost:8080", storage.GaugeMetricName, name, fmt.Sprintf("%f", float64(value)))
		err = createPostRequest(urlGauge)
		if err != nil {
			return err
		}
	}
	return nil
}
