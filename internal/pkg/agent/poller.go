package agent

import (
	"context"
	"fmt"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent/config"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/sirupsen/logrus"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

func RunUpdateMemStatMetrics(ctx context.Context, c *config.AgentConfig, s storage.Store) {
	pollerInterval := time.Duration(c.PollInterval) * time.Second //тесты не проходят с duration
	pollerTicker := time.NewTicker(pollerInterval)
	defer pollerTicker.Stop()

	storeContext, storeCancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer storeCancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollerTicker.C:
			if err := UpdateMemStatMetrics(storeContext, s); err != nil {
				logrus.Errorf("Error update mem stat metrics %v", err)
			}
		}
	}
}

func RunUpdateVirtualMetrics(ctx context.Context, c *config.AgentConfig, s storage.Store) {
	pollerInterval := time.Duration(c.PollInterval) * time.Second //тесты не проходят с duration
	pollerTicker := time.NewTicker(pollerInterval)
	defer pollerTicker.Stop()

	storeContext, storeCancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer storeCancel()

	for {
		select {
		case <-ctx.Done():
			return
		case <-pollerTicker.C:
			if err := UpdateVirtualMetrics(storeContext, s); err != nil {
				logrus.Errorf("Error update virtual metrics %v", err)
			}
		}
	}
}

func UpdateVirtualMetrics(ctx context.Context, m storage.Store) error {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	cpuMetric, err := cpu.Percent(time.Duration(1)*time.Second, true)
	if err != nil {
		return err
	}

	gaugeMetricsMap := make(map[string]metrics.Gauge)
	gaugeMetricsMap["TotalMemory"] = metrics.Gauge(vm.Total)
	gaugeMetricsMap["FreeMemory"] = metrics.Gauge(vm.Free)
	gaugeMetricsMap["CPUutilization1"] = metrics.Gauge(cpuMetric[0])

	var errorsSlice []string

	for k, v := range gaugeMetricsMap {
		err = m.UpdateGaugeMetric(ctx, k, v)
		if err != nil {
			errorsSlice = append(errorsSlice, err.Error())
		}
	}

	if len(errorsSlice) > 0 {
		return fmt.Errorf(strings.Join(errorsSlice, "\n"))
	} else {
		return nil
	}
}

func UpdateMemStatMetrics(ctx context.Context, m storage.Store) error {
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
