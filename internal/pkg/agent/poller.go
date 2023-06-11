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
	var errorsSlice []string

	err := m.UpdateGaugeMetric(ctx, "Alloc", metrics.Gauge(metricsStats.Alloc))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "BuckHashSys", metrics.Gauge(metricsStats.BuckHashSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "BuckHashSys", metrics.Gauge(metricsStats.BuckHashSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "Frees", metrics.Gauge(metricsStats.Frees))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "GCCPUFraction", metrics.Gauge(metricsStats.GCCPUFraction))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "GCSys", metrics.Gauge(metricsStats.GCSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "HeapAlloc", metrics.Gauge(metricsStats.HeapAlloc))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "HeapIdle", metrics.Gauge(metricsStats.HeapIdle))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "HeapInuse", metrics.Gauge(metricsStats.HeapInuse))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "HeapObjects", metrics.Gauge(metricsStats.HeapObjects))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "HeapReleased", metrics.Gauge(metricsStats.HeapReleased))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "HeapSys", metrics.Gauge(metricsStats.HeapSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "LastGC", metrics.Gauge(metricsStats.LastGC))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "Lookups", metrics.Gauge(metricsStats.Lookups))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "MCacheInuse", metrics.Gauge(metricsStats.MCacheInuse))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "MCacheSys", metrics.Gauge(metricsStats.MCacheSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "MSpanInuse", metrics.Gauge(metricsStats.MSpanInuse))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "MSpanSys", metrics.Gauge(metricsStats.MSpanSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "Mallocs", metrics.Gauge(metricsStats.Mallocs))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "NextGC", metrics.Gauge(metricsStats.NextGC))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "NumForcedGC", metrics.Gauge(metricsStats.NumForcedGC))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "NumGC", metrics.Gauge(metricsStats.NumGC))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "OtherSys", metrics.Gauge(metricsStats.OtherSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "PauseTotalNs", metrics.Gauge(metricsStats.PauseTotalNs))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "StackInuse", metrics.Gauge(metricsStats.StackInuse))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "StackSys", metrics.Gauge(metricsStats.StackSys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "Sys", metrics.Gauge(metricsStats.Sys))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "TotalAlloc", metrics.Gauge(metricsStats.TotalAlloc))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}
	err = m.UpdateGaugeMetric(ctx, "RandomValue", metrics.Gauge(rand.Float64()))
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}

	err = m.UpdateCounterMetric(ctx, "PollCount", 1)
	if err != nil {
		errorsSlice = append(errorsSlice, err.Error())
	}

	if len(errorsSlice) > 0 {
		return fmt.Errorf(strings.Join(errorsSlice, "\n"))
	} else {
		return nil
	}
}
