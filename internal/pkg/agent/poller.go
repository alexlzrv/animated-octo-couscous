package agent

import (
	"context"
	"fmt"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"math/rand"
	"runtime"
	"strings"
)

func UpdateMetrics(ctx context.Context, m storage.Store) error {
	var metricsStats runtime.MemStats
	runtime.ReadMemStats(&metricsStats)

	gaugeMetricsMap := make(map[string]metrics.Gauge)
	gaugeMetricsMap["Alloc"] = metrics.Gauge(metricsStats.Alloc)
	gaugeMetricsMap["BuckHashSys"] = metrics.Gauge(metricsStats.BuckHashSys)
	gaugeMetricsMap["Frees"] = metrics.Gauge(metricsStats.Frees)
	gaugeMetricsMap["GCCPUFraction"] = metrics.Gauge(metricsStats.GCCPUFraction)
	gaugeMetricsMap["GCSys"] = metrics.Gauge(metricsStats.GCSys)
	gaugeMetricsMap["HeapAlloc"] = metrics.Gauge(metricsStats.HeapAlloc)
	gaugeMetricsMap["HeapIdle"] = metrics.Gauge(metricsStats.HeapIdle)
	gaugeMetricsMap["HeapInuse"] = metrics.Gauge(metricsStats.HeapInuse)
	gaugeMetricsMap["HeapObjects"] = metrics.Gauge(metricsStats.HeapObjects)
	gaugeMetricsMap["HeapReleased"] = metrics.Gauge(metricsStats.HeapReleased)
	gaugeMetricsMap["HeapSys"] = metrics.Gauge(metricsStats.HeapSys)
	gaugeMetricsMap["Lookups"] = metrics.Gauge(metricsStats.Lookups)
	gaugeMetricsMap["MCacheInuse"] = metrics.Gauge(metricsStats.MCacheInuse)
	gaugeMetricsMap["MCacheSys"] = metrics.Gauge(metricsStats.MCacheSys)
	gaugeMetricsMap["MSpanInuse"] = metrics.Gauge(metricsStats.MSpanInuse)
	gaugeMetricsMap["MSpanSys"] = metrics.Gauge(metricsStats.MSpanSys)
	gaugeMetricsMap["Mallocs"] = metrics.Gauge(metricsStats.Mallocs)
	gaugeMetricsMap["NextGC"] = metrics.Gauge(metricsStats.NextGC)
	gaugeMetricsMap["LastGC"] = metrics.Gauge(metricsStats.LastGC)
	gaugeMetricsMap["NumForcedGC"] = metrics.Gauge(metricsStats.NumForcedGC)
	gaugeMetricsMap["NumGC"] = metrics.Gauge(metricsStats.NumGC)
	gaugeMetricsMap["OtherSys"] = metrics.Gauge(metricsStats.OtherSys)
	gaugeMetricsMap["PauseTotalNs"] = metrics.Gauge(metricsStats.PauseTotalNs)
	gaugeMetricsMap["StackInuse"] = metrics.Gauge(metricsStats.StackInuse)
	gaugeMetricsMap["StackSys"] = metrics.Gauge(metricsStats.StackSys)
	gaugeMetricsMap["Sys"] = metrics.Gauge(metricsStats.Sys)
	gaugeMetricsMap["TotalAlloc"] = metrics.Gauge(metricsStats.TotalAlloc)
	gaugeMetricsMap["RandomValue"] = metrics.Gauge(rand.Float64())

	var errorsSlice []string

	for k, v := range gaugeMetricsMap {
		err := m.UpdateGaugeMetric(ctx, k, v)
		if err != nil {
			errorsSlice = append(errorsSlice, err.Error())
		}
	}

	err := m.UpdateCounterMetric(ctx, "PollCount", 1)
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}

	if len(errorsSlice) > 0 {
		return fmt.Errorf(strings.Join(errorsSlice, "\n"))
	} else {
		return nil
	}
}
