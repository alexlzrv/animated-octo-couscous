package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/agent/config"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/metrics"
	"github.com/mayr0y/animated-octo-couscous.git/internal/pkg/storage"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func RunSendMetric(ctx context.Context, c *config.AgentConfig, s storage.Store) {
	reportInterval := time.Duration(c.ReportInterval) * time.Second //тесты не проходят с duration
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-reportTicker.C:
			for i := 1; i < c.RateLimit; i++ {
				if err := SendMetrics(ctx, s, c.ServerAddress, c.SignKey); err != nil {
					logrus.Errorf("Error send metrics %v", err)
				}
				if err := SendMetricsBatch(ctx, s, c.ServerAddress); err != nil {
					logrus.Errorf("Error send metrics batch %v", err)
				}
				if err := s.ResetCounterMetric(ctx, "PollCount"); err != nil {
					logrus.Errorf("Error reset metrics %v", err)
				}
			}
		}
	}
}

func SendMetrics(ctx context.Context, s storage.Store, serverAddress string, signKey string) error {
	getContext, cancelCtx := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancelCtx()

	url := fmt.Sprintf("http://%s/update/", serverAddress)

	metricMap, err := s.GetMetrics(getContext)
	if err != nil {
		return err
	}

	for _, v := range metricMap {
		err = createPostRequest(url, v, signKey)
		if err != nil {
			return fmt.Errorf("error create post request %v", err)
		}
	}
	return nil
}

func createPostRequest(url string, metric *metrics.Metrics, signKey string) error {
	metric.SetHash(signKey)

	body, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("error encoding metric %v", err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err = gz.Write(body); err != nil {
		return fmt.Errorf("error %s", err)
	}

	gz.Close()

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return fmt.Errorf("error send request %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error client %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("status code not 200")
	}

	return nil
}

func SendMetricsBatch(ctx context.Context, s storage.Store, serverAddress string) error {
	getContext, cancelCtx := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancelCtx()

	metricsMap, err := s.GetMetrics(getContext)
	if err != nil {
		logrus.Errorf("Some error ocured during metrics get: %q", err)
		return err
	}

	metricsBatch := make([]*metrics.Metrics, 0)
	for _, v := range metricsMap {
		metricsBatch = append(metricsBatch, v)
	}

	url := fmt.Sprintf("http://%s/updates/", serverAddress)

	if err = sendBatchJSON(url, metricsBatch); err != nil {
		return fmt.Errorf("error create post request %v", err)
	}
	return nil
}

func sendBatchJSON(url string, metricsBatch []*metrics.Metrics) error {
	body, err := json.Marshal(metricsBatch)
	if err != nil {
		return fmt.Errorf("error encoding metric %v", err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err = gz.Write(body); err != nil {
		return fmt.Errorf("error %s", err)
	}

	gz.Close()

	req, err := http.NewRequest(http.MethodPost, url, &buf)
	if err != nil {
		return fmt.Errorf("error send request %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error client %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %v and not 200", resp.StatusCode)
	}

	return nil
}
